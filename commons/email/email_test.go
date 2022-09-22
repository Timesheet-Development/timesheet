package email

import (
	"fmt"
	"testing"
)

func TestEmail(t *testing.T) {

	var (
		Subject = "Password Authentication"

		Body = "Don't share this to anyone"
	)

	email := Email{
		EmailFrom: EmailFromVar,
		Password:  PasswordVar,
		EmailTo:   EmailToVar,
		Subject:   Subject,
		Body:      Body,
		Host:      HOST,
		Port:      PORT,
		Address:   HOST + ":" + PORT,
	}

	err := SendEmail(&email)
	if err != nil {
		t.Errorf("error while performing send email %v", err)
	}
	fmt.Println("Email Sent Successfully")
}
