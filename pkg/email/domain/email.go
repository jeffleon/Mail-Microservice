package domain

type EMail struct {
	Domain      string
	FromAddress string
	FromName    string
}

type Message struct {
	From     string `json:"from"`
	FromName string `json:"from_name"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

type MailerRepository interface {
	SendEmail(msg Message) error
}
