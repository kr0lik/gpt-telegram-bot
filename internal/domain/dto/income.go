package dto

type MessageId string
type ChatId string

type Income struct {
	MessageId MessageId
	ChatId    ChatId
	UserId    string
	Message   string
	Command   string
	Callback  IncomeCallback
	ImagePath string
	Caption   string
	AudioPath string
}

type IncomeCallback struct {
	Id        string
	MessageId MessageId
	Command   string
}
