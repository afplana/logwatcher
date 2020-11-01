package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
	logger "github.com/sirupsen/logrus"
)

// Smtp server constants
const smtpServer string = "" // smtp server addr
const smtpPort string = "25"
const servername = "/etc/hostname"

func getHostname() string {
	if hostname, err := ioutil.ReadFile(servername); err == nil {
		return string(hostname)
	}
	return ""
}

// Send email notification
func emailNotify(str string) error {
	replacer := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")
	msgBody := []byte(str) // Message of notification

	// email details
	sender := "no-reply@bosch.io"
	to := strings.Split("", ",") // comma separated string with receivers email addr
	hostname := replacer.Replace(getHostname())
	subject := "Errors in server " + hostname //

	// Dial smtp server for retrieve a connection
	conn, dialErr := smtp.Dial(smtpServer + ":" + smtpPort)
	if dialErr != nil {
		return dialErr
	}
	defer conn.Close()

	// Set sender
	conn.Mail(replacer.Replace(sender))

	// Build receipient list
	for i := range to {
		if recpErr := conn.Rcpt(replacer.Replace(to[i])); recpErr != nil {
			logger.Info("Error setting receipient " + to[i])
		}
	}

	wr, errData := conn.Data()
	if errData != nil {
		return errData
	}

	// Message structure
	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + sender + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString(msgBody)

	// Write message into email body
	wr.Write([]byte(msg))

	// Close connection
	if err := wr.Close(); err != nil {
		return err
	}

	conn.Quit()
	return nil
}

// Watch logs files for errors
func watchLogs(path string) {
	errorPattern := regexp.MustCompile("^.*\\[ERROR\\].*") // matches with error in psing logs files

	// Read from tail of file
	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		logger.Error(err)
	}

	// if any line is error send a email notification
	for line := range t.Lines {
		if errorPattern.MatchString(line.Text) {
			if er := emailNotify(line.Text); er != nil {
				logger.Info("Email notification was not sent. Cause: ", er)
			} else {
				logger.Info("Notification sent to recipient list!")
			}
		}
	}
}

func main() {
	logsPath := os.Args[1]
	if logsPath != "" {
		watchLogs(logsPath)
	}
}
