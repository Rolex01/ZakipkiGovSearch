package mail

import (
	"gopkg.in/gomail.v2"
)

var d *gomail.Dialer

func Init(host string, port int, username, password string) {
	d = gomail.NewDialer(host, port, username, password)
}

func SendMessage(to, subject, message string) error {

	m := gomail.NewMessage()
	
	m.SetHeader("From", d.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	return d.DialAndSend(m)
}
