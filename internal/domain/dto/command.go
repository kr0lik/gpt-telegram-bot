package dto

import "gpt-telegran-bot/internal/domain/enum"

type Command struct {
	Id            string
	Description   string
	Sub           []Command
	IsInstruction bool
}

var (
	HelpCommand       = Command{Id: enum.CommandHelp, Description: "Show available commands"}
	StatusCommand     = Command{Id: enum.CommandStatus, Description: "Show info"}
	ChatCommand       = Command{Id: enum.CommandChat, Description: "Turn to chat generation model", Sub: []Command{ChatNewCommand}}
	ChatNewCommand    = Command{Id: enum.CommandChatNew, Description: "Clear conversation history"}
	TextCommand       = Command{Id: enum.CommandText, Description: "Turn to text generation model"}
	TextEditCommand   = Command{Id: enum.CommandTextEdit, Description: "Turn to edit text generation model"}
	CodeEditCommand   = Command{Id: enum.CommandCodeEdit, Description: "Turn to edit code generation model"}
	ImageCommand      = Command{Id: enum.CommandImage, Description: "Turn to image generation model", Sub: []Command{ImageSizeCommand, ImageCountCommand}}
	ImageEditCommand  = Command{Id: enum.CommandImageEdit, Description: "Turn to edit image generation model", Sub: []Command{ImageSizeCommand, ImageCountCommand}}
	ImageSizeCommand  = Command{Id: enum.CommandImageSize + " n", Description: "Set image size"}
	ImageCountCommand = Command{Id: enum.CommandImageCount + " n", Description: "Set count images"}
	SpeechCommand     = Command{Id: enum.CommandSpeech, Description: "Turn to audio transcription model"}
)
