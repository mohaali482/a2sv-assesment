package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
)

type EmailService interface {
	Send(to, subject, body string) error
}

type EmailServiceImpl struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(host string, port int, username, password string) EmailService {
	return &EmailServiceImpl{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (e *EmailServiceImpl) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	fromMsg := fmt.Sprintf("From: <%s>\r\n", e.Username)
	toMSG := fmt.Sprintf("To: <%s>\r\n", to)
	subject = fmt.Sprintf("Subject: %s\r\n", subject)

	message := []byte(fromMsg + toMSG + subject + "\r\n" + body)

	addr := fmt.Sprintf("%s:%s", e.Host, strconv.Itoa(e.Port))
	err := smtp.SendMail(addr, auth, e.Username, []string{to}, message)
	if err != nil {
		log.Default().Println("Error while sending email", err)
		return errors.New("failed to send an email")
	}
	return nil

}
