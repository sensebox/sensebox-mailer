package main

import (
	"html/template"
	"strconv"
	"time"

	"github.com/honeybadger-io/honeybadger-go"
	// should be "github.com/jordan-wright/email"
	// but we wait until https://github.com/jordan-wright/email/pull/61 is merged
	"github.com/lovego/email"
)

var (
	branch, ts, hash string
)

type senseBoxMailerServer struct {
	Daemon chan *email.Email
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func logStartup() {
	LogInfo("sensebox-mailer")
	var isoTime time.Time

	iTs, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		LogInfo("startup", "version:", branch, "??? (", ts, ")", hash)
		return
	}
	isoTime = time.Unix(iTs, 0)

	//                             branch, timestamp, hash
	LogInfo("startup", "version:", branch, isoTime.Format(time.RFC3339), hash)
}

func main() {
	defer honeybadger.Monitor()

	logStartup()

	initConfigFromEnv()
	loadTranslationsJson()
	templates.Option("missingkey=error")

	mailer := senseBoxMailerServer{}
	mailer.startMailerDaemon()
	mailer.StartHTTPSServer()
	defer close(mailer.Daemon)
}
