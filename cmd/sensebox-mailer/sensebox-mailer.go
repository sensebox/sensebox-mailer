package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/honeybadger-io/honeybadger-go"
	"github.com/sensebox/sensebox-mailer/mailer"
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

func fetchLatestTemplatesFromGithub() {
	ticker := time.NewTicker(60000 * time.Millisecond)

	for range ticker.C {
		cmd := exec.Command("git", "pull", "origin", "main")
		cmd.Dir = "./sensebox-mailer-templates"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error pulling git repository")
		}
	}
}

func main() {
	defer honeybadger.Monitor()

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

	// Check if templates folder exists
	if _, err := os.Stat("./sensebox-mailer-templates"); os.IsNotExist(err) {
		fmt.Println("Templates do not exists; Go and clone repository")
		cmd := exec.Command("git", "clone", "git@github.com:sensebox/sensebox-mailer-templates.git")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Git clone failed :(")
		} else {
			fmt.Println("Successfully cloned templates")
		}
	} else {
		fmt.Println("Templates exists; go and start pulling updates")
	}

	// Start routine to fetch latest templates
	go fetchLatestTemplatesFromGithub()

	err := mailer.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
