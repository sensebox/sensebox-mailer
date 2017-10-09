package main

import (
	"html/template"

	"github.com/honeybadger-io/honeybadger-go"
	"gopkg.in/gomail.v2"
)

type senseBoxMailerServer struct {
	Daemon chan *gomail.Message
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	defer honeybadger.Monitor()
	initConfigFromEnv()
	loadTranslationsJson()
	templates.Option("missingkey=error")

	mailer := senseBoxMailerServer{}
	mailer.startMailerDaemon()
	mailer.StartHTTPSServer()
	defer close(mailer.Daemon)
}
