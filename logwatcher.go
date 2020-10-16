package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
)

// Smtp server constants
const smtpServer string = "" // smtp server addr
const smtpPort string = "25"
const servername = "/etc/hostname"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getHostname() string {
	hostname, err := ioutil.ReadFile(servername)
	checkError(err)
	return string(hostname)
}

// Send email notification
func emailNotify(str string) {
	replacer := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")
	msgBody := []byte(str) // Message of notification

	// email details
	sender := "no-reply@bosch.io"
	toList := strings.Split("", ",") // comma separated string with receivers email addr
	hostname := replacer.Replace(getHostname())
	subject := " " + hostname //

	// Dial smtp server for retrieve a connection
	conn, dialErr := smtp.Dial(smtpServer + ":" + smtpPort)
	checkError(dialErr)
	defer conn.Close()

	// Set sender
	senderErr := conn.Mail(replacer.Replace(sender))
	checkError(senderErr)

	// Build receipient list
	for i := range toList {
		toList[i] = replacer.Replace(toList[i])
		recpErr := conn.Rcpt(toList[i])
		checkError(recpErr)
	}

	wr, bodyErr := conn.Data()
	checkError(bodyErr)

	// Message structure
	msg := "To: " + strings.Join(toList, ",") + "\r\n" +
		"From: " + sender + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString(msgBody)

	// Write message into email body
	_, writeErr := wr.Write([]byte(msg))
	checkError(writeErr)

	// Close connection
	closeWErr := wr.Close()
	checkError(closeWErr)
	conn.Quit()
}

// Watch logs files for errors
func watchLogs(path string) {
	errorPattern := regexp.MustCompile("^.*\\[ERROR\\].*") // matches with error in psing logs files

	// Read from tail of file
	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true})

	checkError(err)
	// if any line is error send a email notification
	for line := range t.Lines {
		if errorPattern.MatchString(line.Text) {
			emailNotify(line.Text)
		}
	}
}

func main() {
	logsPath := os.Args[1]
	if logsPath != "" {
		watchLogs(logsPath)
	}
}
