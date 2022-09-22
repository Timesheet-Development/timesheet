package email

import gomail "gopkg.in/gomail.v2"

const (
	HOST = "smtp.gmail.com"

	PORT = "587"
)

var (
	EmailFromVar = "jayanthsravan00@gmail.com"
	PasswordVar  = "Kalyani74$"
	EmailToVar   = []string{"sravannemala@gmail.com"}
)

type Email struct {
	EmailFrom, Password, Subject, Body, Host, Port, Address string
	EmailTo                                                 []string
}

func SendEmail(email *Email) error {
	var err error
	msg := gomail.NewMessage()
	msg.SetHeader("From", email.EmailFrom)
	msg.SetHeader("To", email.EmailTo...)
	msg.SetHeader("Subject", email.Subject)
	msg.SetHeader("text/html", email.Body)

	n := gomail.NewDialer(email.Host, 587, email.EmailFrom, email.Password)
	err = n.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
