package mail

/*
	Modified code based on article https://hackernoon.com/golang-sendmail-sending-mail-through-net-smtp-package-5cadbe2670e0
*/

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

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
	msg += MIME + "\r\n"
	msg += "\r\n" + mail.body

	return msg
}

func (m *Mail) parseTemplate(file string, msg interface{}) error {
	filePath, _ := filepath.Abs("./mail/templates/" + file)
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}
	bfr := new(bytes.Buffer)
	if err = t.Execute(bfr, msg); err != nil {
		return err
	}
	m.body = bfr.String()
	return nil
}

// Send sends email
func Send(to []string, subj string, msgData interface{}, htmlTemp string) {
	mail := &Mail{}

	err := mail.parseTemplate(htmlTemp, msgData)
	if err != nil {
		log.Fatal(err)
	}

	mail.sender = os.Getenv("mail_user")
	mail.recievers = to
	mail.subject = subj

	// Build email
	finalMail := mail.Message()

	mServer := MailServer{
		host: os.Getenv("mail_server"),
		port: os.Getenv("mail_port"),
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
