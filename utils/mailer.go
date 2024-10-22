package utils

import (
	"bytes"
	"embed"
	"github.com/go-mail/mail/v2"
	"html/template"
	"log"
	"time"
)

//go:embed template/*
var templateFS embed.FS

type Mailer struct {
	Dialer *mail.Dialer
	Sender string
}

func NewMailer(host, username, password, sender string, port int) Mailer {
	dialer := mail.NewDialer(host, port, username, password)

	return Mailer{
		Dialer: dialer,
		Sender: sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "template/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.Sender) // Ensure this is a valid email address
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for i := 0; i < 5; i++ {
		err = m.Dialer.DialAndSend(msg)
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	log.Println(err)
	return err
}
