package mail

/*
	Modified code based on article https://hackernoon.com/golang-sendmail-sending-mail-through-net-smtp-package-5cadbe2670e0
*/

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {
	sender    string
	recievers []string
	subject   string
	body      string
}

type MailServer struct {
	host string
	port string
}

func (m *MailServer) Server() string {
	return m.host + ":" + m.port
}

func (mail *Mail) Message() string {
	msg := ""
	msg += fmt.Sprintf("From: %s\r\n", mail.sender)

	if len(mail.recievers) > 0 {
		msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.recievers, ";"))
	}

	msg += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	msg += "\r\n" + mail.body

	return msg
}

// Send sends email
func Send(to []string, subj string, msg string) {
	mail := Mail{
		sender:    os.Getenv("mail_user"),
		recievers: to,
		subject:   subj,
		body:      msg,
	}

	// Build email
	finalMail := mail.Message()

	mServer := MailServer{
		host: "smtp.gmail.com",
		port: "465",
	}

	// New smtp auth that uses password from .env
	smtpAuth := smtp.PlainAuth("", mail.sender, os.Getenv("mail_password"), mServer.host)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mServer.host,
	}

	// New tcp connection
	conn, err := tls.Dial("tcp", mServer.Server(), tlsConfig)
	if err != nil {
		log.Panic(err)
	}

	// New
	client, err := smtp.NewClient(conn, mServer.host)
	if err != nil {
		log.Panic(err)
	}

	// Authenticate
	if err = client.Auth(smtpAuth); err != nil {
		log.Panic(err)
	}

	if err = client.Mail(mail.sender); err != nil {
		log.Panic(err)
	}

	for _, k := range mail.recievers {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(finalMail))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Email sent successfully")

}
