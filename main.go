package main

import (
	"log"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Send email notification
func emailNotify(str string) {
	msg := []byte(str) // Message of notification

	// email details
	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_SENDER_PASSWORD")
	toList := strings.Split(os.Getenv("EMAIL_RECEIVER_LIST"), ",")

	// smtp server config.
	host := os.Getenv("EMAIL_CONFIG_HOST") // for testing should use "smtp.gmail.com"
	port := os.Getenv("EMAIL_CONFIG_PORT") // for testing 587

	// authentication details.
	auth := smtp.PlainAuth("", sender, password, host)

	// Sending email.
	err := smtp.SendMail(host+":"+port, auth, sender, toList, msg)
	checkError(err)

	log.Println("Email Sent Successfully!")
}

// Watch logs files for errors
func watchLogs(path string) {
	errorPattern := regexp.MustCompile("^.*\\[error\\].*")

	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true})
	checkError(err)

	for line := range t.Lines {
		if errorPattern.MatchString(line.Text) {
			emailNotify(line.Text)
		}
	}
}

func main() {
	logsPath := os.Args[1]
	watchLogs(logsPath)
}
