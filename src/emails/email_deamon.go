package emails

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
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
	server.Host = viper.GetString("email.host")
	server.Port = viper.GetInt("email.port")
	server.Username = viper.GetString("email.username")
	server.Password = viper.GetString("email.password")
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

	email_subject := "Recuperacion de contraseña - CarMind"

	email_body := "Hola, tu token de recuperacion de contraseña es: <p><h3>{{pass}}</h3></p>"

	for {
		select {
		case to, ok := <-ch:
			if !ok {
				return
			}
			email := ed._createMessage(email_from, email_subject, email_body, to.Token, nil)
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

	email.SetBody(mail.TextHTML, strings.Replace(body, "{{pass}}", token, -1))

	// add inline
	if file != nil {
		email.Attach(file)
	}
	return email
}
