package messenger

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gpt-telegran-bot/internal/domain/dto"
	"gpt-telegran-bot/internal/domain/service"
	"gpt-telegran-bot/internal/infrastructure/util"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	telegramUpdatesTimeout   = 60
	telegramMaxPhotoSize     = 4 * 1024 * 1024
	telegramMaxMessageLength = 4096

	telegramParseModeError = "Can't find end of the entity starting at byte offset"
)

type Telegram struct {
	api          *tgbotapi.BotAPI
	downloadPath string
	allowedUsers []string
}

type TelegramConfig struct {
	ApiToken     string
	DownloadPath string
	AllowedUsers []string
}

func NewTelegram(config *TelegramConfig) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initilize Telegram: %w", err)
	}

	return &Telegram{
		api:          bot,
		downloadPath: config.DownloadPath,
		allowedUsers: config.AllowedUsers,
	}, nil
}

func (t *Telegram) PrintCommands(title string, printCommands [][]dto.Command, chatId dto.ChatId, commands [][]dto.Command) {
	message := title

	for _, row := range printCommands {
		message += "\n------------------"
		for _, cmd := range row {
			message += fmt.Sprintf("\n/%s - %s", cmd.Id, cmd.Description)
			for _, sub := range cmd.Sub {
				message += fmt.Sprintf("\n---- /%s - %s", sub.Id, sub.Description)
			}
		}
	}

	t.Send(message, chatId, commands)
}

func (t *Telegram) Listen(ctx context.Context) (service.UpdatesChannel, error) {
	ch := make(chan dto.Income)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
			}

			u := tgbotapi.NewUpdate(0)
			u.Timeout = telegramUpdatesTimeout

			updates, err := t.api.GetUpdatesChan(u)
			if err != nil {
				log.Printf("failed get updates from telegram: %v", err)
				close(ch)
			}

			for update := range updates {
				t.processUpdate(ch, &update)
			}
		}
	}()

	return ch, nil
}

func (t *Telegram) processUpdate(ch chan<- dto.Income, update *tgbotapi.Update) {
	if update.CallbackQuery != nil {
		callbackData := strings.Split(update.CallbackQuery.Data, " ")

		ch <- dto.Income{
			MessageId: t.messageIdFromTelegram(update.CallbackQuery.Message.MessageID),
			ChatId:    t.chatIdFromTelegram(update.CallbackQuery.Message.Chat.ID),
			UserId:    update.CallbackQuery.Message.ReplyToMessage.Chat.UserName,
			Callback: dto.IncomeCallback{
				Id:        update.CallbackQuery.ID,
				MessageId: dto.MessageId(callbackData[1]),
				Command:   callbackData[0],
			},
		}
	}

	if update.Message == nil {
		return
	}

	if !t.isAllowedUser(update.Message) {
		return
	}

	ch <- dto.Income{
		MessageId: t.messageIdFromTelegram(update.Message.MessageID),
		ChatId:    t.chatIdFromTelegram(update.Message.Chat.ID),
		UserId:    update.Message.Chat.UserName,
		Message:   update.Message.Text,
		Command:   update.Message.Command(),
		Callback:  dto.IncomeCallback{},
		ImagePath: t.getImage(update),
		Caption:   update.Message.Caption,
		AudioPath: t.getAudio(update),
	}
}

func (t *Telegram) isAllowedUser(message *tgbotapi.Message) bool {
	for _, allowedUser := range t.allowedUsers {
		if allowedUser == message.From.UserName {
			return true
		}
	}

	t.Replay(message.From.UserName+" not allowed", t.messageIdFromTelegram(message.MessageID), t.chatIdFromTelegram(message.Chat.ID), [][]dto.Command{})

	return false
}

func (t *Telegram) Send(message string, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId {
	msg := tgbotapi.NewMessage(t.chatIdToTelegram(chatId), message)

	if len(commands) > 0 {
		msg.ReplyMarkup = t.getKeyboardMarkup(commands)
	}

	result, err := t.api.Send(msg)
	if err != nil {
		log.Printf("failed to send message: %v", err)
		log.Printf("failed to send message text: %s", message)

		return ""
	}

	return t.messageIdFromTelegram(result.MessageID)
}

func (t *Telegram) Replay(message string, replayId dto.MessageId, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId {
	msg := tgbotapi.NewMessage(t.chatIdToTelegram(chatId), message)
	msg.ReplyToMessageID = t.messageIdToTelegram(replayId)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if len(commands) > 0 {
		msg.ReplyMarkup = t.getKeyboardMarkup(commands)
	}

	result, err := t.api.Send(msg)
	if err != nil {
		if strings.Contains(err.Error(), telegramParseModeError) {
			msg.ParseMode = ""
			result, err := t.api.Send(msg)
			if err != nil {
				return t.Send(message, chatId, commands)
			}

			return t.messageIdFromTelegram(result.MessageID)
		}

		return t.Send(message, chatId, commands)
	}

	return t.messageIdFromTelegram(result.MessageID)
}

func (t *Telegram) Replace(messageId dto.MessageId, newMessage string, replayId dto.MessageId, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId {
	t.delete(messageId, chatId)
	return t.Replay(newMessage, replayId, chatId, commands)
}

func (t *Telegram) ReplaceWithPhotos(messageId dto.MessageId, urls []string, replayId dto.MessageId, chatId dto.ChatId, commands [][]dto.Command) dto.MessageId {
	media := []interface{}{}
	for _, url := range urls {
		media = append(media, tgbotapi.NewInputMediaPhoto(url))
	}

	msg := tgbotapi.NewMediaGroup(t.chatIdToTelegram(chatId), media)
	msg.ReplyToMessageID = t.messageIdToTelegram(replayId)

	if len(commands) > 0 {
		msg.ReplyMarkup = t.getKeyboardMarkup(commands)
	}

	result, err := t.api.Send(msg)

	t.delete(messageId, chatId)

	if err != nil {
		return t.Replay(err.Error(), replayId, chatId, commands)
	}

	return t.messageIdFromTelegram(result.MessageID)
}

func (t *Telegram) StartEdit(message string, replayId dto.MessageId, chatId dto.ChatId, callbacks [][]dto.Callback, commands [][]dto.Command) dto.MessageId {
	msg := tgbotapi.NewMessage(t.chatIdToTelegram(chatId), message)
	msg.ReplyToMessageID = t.messageIdToTelegram(replayId)

	if len(callbacks) > 0 {
		msg.ReplyMarkup = *t.getInlineKeyboardMarkup(callbacks)
	}

	result, err := t.api.Send(msg)
	if err != nil {
		return t.Send(message, chatId, commands)
	}

	return t.messageIdFromTelegram(result.MessageID)
}

func (t *Telegram) Edit(messageId dto.MessageId, newMessage string, replayId dto.MessageId, chatId dto.ChatId, callbacks [][]dto.Callback, commands [][]dto.Command) dto.MessageId {
	if len(newMessage) > telegramMaxMessageLength {
		runes := []rune(newMessage)
		numRunes := utf8.RuneCountInString(newMessage)

		prevIndex := numRunes - telegramMaxMessageLength*2
		if prevIndex < 0 {
			prevIndex = 0
		}

		prev := string(runes[prevIndex:])
		replayId := t.Replace(messageId, prev, replayId, chatId, commands)

		lastIndex := numRunes - telegramMaxMessageLength
		lastText := string(runes[lastIndex:])
		return t.StartEdit(lastText, replayId, chatId, callbacks, commands)
	}

	msg := tgbotapi.NewEditMessageText(t.chatIdToTelegram(chatId), t.messageIdToTelegram(messageId), newMessage)

	if len(callbacks) > 0 {
		msg.ReplyMarkup = t.getInlineKeyboardMarkup(callbacks)
	}

	if _, err := t.api.Send(msg); err != nil {
		newMessageId := t.StartEdit(newMessage, replayId, chatId, callbacks, commands)
		t.delete(messageId, chatId)
		return newMessageId
	}

	return messageId
}

func (t *Telegram) Callback(callbackId, message string) {
	answer := tgbotapi.NewCallback(callbackId, message)
	_, err := t.api.AnswerCallbackQuery(answer)
	if err != nil {
		log.Printf("failed to send callback message: %v", err)
	}
}

func (t *Telegram) delete(messageId dto.MessageId, chatId dto.ChatId) {
	msg := tgbotapi.NewDeleteMessage(t.chatIdToTelegram(chatId), t.messageIdToTelegram(messageId))
	if _, err := t.api.Send(msg); err != nil {
		log.Printf("failed to delete message: %v", err)
	}
}

func (t *Telegram) chatIdToTelegram(chatId dto.ChatId) int64 {
	res, err := strconv.ParseInt(string(chatId), 10, 64)
	if err != nil {
		log.Fatalf("failed to convert chat id %v: %v", chatId, err)
	}

	return res
}

func (t *Telegram) chatIdFromTelegram(chatId int64) dto.ChatId {
	return dto.ChatId(strconv.FormatInt(chatId, 10))
}

func (t *Telegram) messageIdToTelegram(messageId dto.MessageId) int {
	res, err := strconv.Atoi(string(messageId))
	if err != nil {
		log.Printf("failed to convert message id %v: %v", messageId, err)
	}

	return res
}

func (t *Telegram) messageIdFromTelegram(messageId int) dto.MessageId {
	return dto.MessageId(strconv.Itoa(messageId))
}

func (t *Telegram) getKeyboardMarkup(commands [][]dto.Command) tgbotapi.ReplyKeyboardMarkup {
	buttons := make([][]tgbotapi.KeyboardButton, 0)

	for _, row := range commands {
		var buttonRow []tgbotapi.KeyboardButton

		for _, cmd := range row {
			text := cmd.Id

			if !cmd.IsInstruction {
				text = "/" + text
			}

			buttonRow = append(buttonRow, tgbotapi.NewKeyboardButton(text))
		}

		buttons = append(buttons, buttonRow)
	}

	return tgbotapi.NewReplyKeyboard(buttons...)
}

func (t *Telegram) getInlineKeyboardMarkup(callbacks [][]dto.Callback) *tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, row := range callbacks {
		var buttonRow []tgbotapi.InlineKeyboardButton

		for _, clb := range row {
			fn := fmt.Sprintf("%s %s", clb.Id, clb.MessageId)
			buttonRow = append(buttonRow, tgbotapi.NewInlineKeyboardButtonData(clb.Description, fn))
		}

		buttons = append(buttons, buttonRow)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	return &keyboard
}

func (t *Telegram) getImage(update *tgbotapi.Update) string {
	if update.Message.Photo == nil && update.Message.Document == nil {
		return ""
	}

	var fileId string

	if update.Message.Document != nil && strings.Contains(update.Message.Document.MimeType, "image") {
		fileId = update.Message.Document.FileID
	}

	if update.Message.Photo != nil {
		maxSize := 0

		for _, photo := range *update.Message.Photo {
			if photo.FileSize > telegramMaxPhotoSize || photo.FileSize < maxSize {
				continue
			}

			maxSize = photo.FileSize
			fileId = photo.FileID
		}
	}

	if fileId == "" {
		return ""
	}

	path, err := t.downloadFile(update.Message.Chat.ID, fileId)
	if err != nil {
		log.Printf("error while download photo file: %v", err)
	}

	return path
}

func (t *Telegram) getAudio(update *tgbotapi.Update) string {
	if update.Message.Audio == nil && update.Message.Voice == nil {
		return ""
	}

	var fileId string

	if update.Message.Audio != nil {
		fileId = update.Message.Audio.FileID
	}

	if update.Message.Voice != nil {
		fileId = update.Message.Voice.FileID
	}

	if fileId == "" {
		return ""
	}

	path, err := t.downloadFile(update.Message.Chat.ID, fileId)
	if err != nil {
		log.Printf("error while download audio file: %v", err)
	}

	return path
}

func (t *Telegram) downloadFile(chatId int64, fileId string) (string, error) {
	fileConfig := tgbotapi.FileConfig{FileID: fileId}
	file, err := t.api.GetFile(fileConfig)
	if err != nil {
		log.Printf("error while getting photo file: %v", err)
	}

	folderPath := strings.TrimRight(t.downloadPath, "/") + "/" + string(t.chatIdFromTelegram(chatId))
	filePath := folderPath + "/" + filepath.Base(file.FilePath)

	if err := util.DownloadFile(file.Link(t.api.Token), filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
