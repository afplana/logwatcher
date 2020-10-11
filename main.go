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
	Sender := os.Getenv("EMAIL_SENDER")
	Pass := os.Getenv("EMAIL_SENDER_PASSWORD")
	ToList := strings.Split(os.Getenv("EMAIL_RECEIVER_LIST"), ",")

	// smtp server config.
	Host := os.Getenv("EMAIL_CONFIG_HOST") // for testing should use "smtp.gmail.com"
	Port := os.Getenv("EMAIL_CONFIG_PORT") // for testing 587

	// authentication details.
	auth := smtp.PlainAuth("", Sender, Pass, Host)

	// Sending email.
	err := smtp.SendMail(Host+":"+Port, auth, Sender, ToList, msg)
	checkError(err)

	log.Println("Email Sent Successfully!")
}

// Watch logs files for errors
func WatchLogs(path string) {
	ErrorPattern := regexp.MustCompile("^.*\\[error\\].*")

	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true})
	checkError(err)

	for line := range t.Lines {
		if ErrorPattern.MatchString(line.Text) {
			emailNotify(line.Text)
		}
	}
}

func main() {
	WatchLogs("/home/aklan/VSProjects/tutor")
}
