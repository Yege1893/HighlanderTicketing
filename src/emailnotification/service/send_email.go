package service

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func SendEmail( /*toList []string, subject string, body string*/ ) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	from := mail.Address{
		Name:    "Highlander Ticketing",
		Address: os.Getenv("EMAIL_ADRESS"),
	}

	fmt.Println(from)

	toList := []string{"yannick.ege@web.de"}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = strings.Join(toList, ", ")
	header["Subject"] = "subject"

	body := "test"

	message := ""
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// über os variablen holen
	smtpServer := "smtp.web.de"
	smtpPort := "587"
	password := os.Getenv("EMAIL_PW")

	auth := smtp.PlainAuth("", from.Address, password, smtpServer)

	err1 := smtp.SendMail(smtpServer+":"+smtpPort, auth, from.Address, toList, []byte(message))
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	fmt.Println("E-Mail erfolgreich gesendet.")
}
