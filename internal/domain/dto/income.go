package dto

type MessageId string
type ChatId string

type Income struct {
	MessageId MessageId
	ChatId    ChatId
	UserId    string
	Message   string
	Command   string
	ImagePath string
	Caption   string
	AudioPath string
}
