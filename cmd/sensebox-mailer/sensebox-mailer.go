package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sensebox/sensebox-mailer/mailer"
	"github.com/sensebox/sensebox-mailer/mailer/templates"
	// should be "github.com/jordan-wright/email"
	// but we wait until https://github.com/jordan-wright/email/pull/61 is merged
)

var (
	branch, ts, hash string
)

func logStartup() {
	var timestamp string

	iTs, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		timestamp = fmt.Sprintf("??? (%s)", ts)
	}
	timestamp = time.Unix(iTs, 0).UTC().Format(time.RFC3339)

	fmt.Printf("sensebox-mailer startup. Version: %s %s %s\n", branch, timestamp, hash)
}

func main() {
	logStartup()
	caCert, serverCert, serverKey, smtpServer, smtpUser, smtpPassword, fromDomain, smtpPort, errors := initConfigFromEnv()
	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	mailer := mailer.MailerServer{
		CaCert:       caCert,
		ServerCert:   serverCert,
		ServerKey:    serverKey,
		SMTPServer:   smtpServer,
		SMTPPort:     smtpPort,
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
		FromDomain:   fromDomain,
	}

	err := templates.CloneTemplatesFromGitHub()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start routine to fetch latest templates
	go templates.FetchLatestTemplatesFromGithub()

	err = mailer.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
