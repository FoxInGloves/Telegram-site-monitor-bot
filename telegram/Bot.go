package telegram

type Bot interface {
	SendMessage(text string) error
	Updates() <-chan Update
}

type Update struct {
	Text    string
	Command string
}
