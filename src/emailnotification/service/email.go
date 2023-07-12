package service

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gitlab.reutlingen-university.de/ege/highlander-ticketing-go-ss2023/src/emailnotification/model"
)

func CreateEmail(emailContenct model.EmialContent, subject string) (string, string, string) {
	if subject == "confirm" {
		return emailContenct.Emailadress, fmt.Sprintf("Hallo Herr/Frau, %s\r\nHiermit bestaetigen wird deine Bestellung fuer das VFB Spiel", emailContenct.Name), "Confirm Order"
	}
	if subject == "cancel" {
		return emailContenct.Emailadress, fmt.Sprintf("Hallo Herr/Frau, %s\r\nHiermit bestaetigen wird die Stornierung deiner Bestellung fuer das VFB Spiel in %s, am %s", emailContenct.Name, emailContenct.Location, emailContenct.Date), "Confirm Cancelation"
	}
	return "", "", ""
}
func SendEmail(receiver string, body string, subject string) error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		return fmt.Errorf("Error loading .env file")
	}

	from := mail.Address{
		Name:    "Highlander Ticketing",
		Address: os.Getenv("EMAIL_ADRESS"),
	}

	toList := []string{}
	toList = append(toList, receiver)

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = strings.Join(toList, ", ")
	header["Subject"] = subject

	message := ""
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	smtpServer := "smtp.web.de"
	smtpPort := "587"
	password := os.Getenv("EMAIL_PW")

	auth := smtp.PlainAuth("", from.Address, password, smtpServer)

	err1 := smtp.SendMail(smtpServer+":"+smtpPort, auth, from.Address, toList, []byte(message))
	if err1 != nil {
		return err1
	}
	return nil
}
