package emails

import (
	"fmt"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

var ChannelEmails chan Recuperacion

type EmailDeamon struct {
}

func (ed *EmailDeamon) Init() {
	ChannelEmails = make(chan Recuperacion)
	go ed.SendEmailChannel(ChannelEmails)
}

func (ed *EmailDeamon) SendEmailChannel(ch <-chan Recuperacion) {

	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = "c2090187.ferozo.com"
	server.Port = 465
	server.Username = "recover@mindiasoft.com"
	server.Password = "Carmind2022"
	server.Encryption = mail.EncryptionSSL

	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		panic(err)
	}

	email_from := "From Example <recover@mindiasoft.com>"

	email_subject := "New Go Email (Prueba del channel)"

	email_body := "Hello <b>Alex</b> and <i>Dan</i>!"

	email_file := &mail.File{FilePath: "C:/Users/Alexis Fonzo/Pictures/Yo-en-la-nao.jpeg", Name: "Yo-en-la-nao.jpeg", Inline: true}

	for {
		select {
		case to, ok := <-ch:
			if !ok {
				return
			}
			email := ed._createMessage(email_from, email_subject, email_body, to.Token, email_file)
			email.AddTo(to.Email)

			err := email.Send(smtpClient)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}
	}
}

func (ed *EmailDeamon) _createMessage(from, subject, body, token string, file *mail.File) *mail.Email {
	email := mail.NewMSG()
	email.SetFrom(from).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, body+"  "+token)

	// add inline
	email.Attach(file)
	return email
}
