package models

//Mail models definition
type Mail struct {
	From, subject, Body, mime string
	To                        []string
}

const (
	mime = "MIME-version: 1.0;\r\nContent-Type: text/html; charset=UTF-8\r\n"
)

//NewMail create a new mail
func NewMail(from string, to []string, subject string) *Mail {
	return &Mail{
		From:    from,
		To:      to,
		Body:    "",
		subject: subject,
		mime:    mime,
	}
}

//GetHeaders returns the formatted headers
func (me *Mail) GetHeaders() string {
	return "From: " + me.From + "\r\nTo: " + me.To[0] + "\r\nSubject: " + me.subject + "\r\n" + mime + "\r\n"
}
