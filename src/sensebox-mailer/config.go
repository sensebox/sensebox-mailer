package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/honeybadger-io/honeybadger-go"
)

const envPrefix = "SENSEBOX_MAILER_"

var ConfigCaCertBytes, ConfigServerCertBytes, ConfigServerKeyBytes []byte
var ConfigSmtpServer, ConfigSmtpUser, ConfigSmtpPassword, ConfigFromDomain string
var ConfigSmtpPort int
var Translations map[string]interface{}

func initConfigFromEnv() {
	errors := make([]error, 0)

	// try to configure honeybadger integration..
	honeybadgerApiKey, _ := getStringFromEnv("HONEYBADGER_APIKEY")
	if honeybadgerApiKey != "" {
		honeybadger.Configure(honeybadger.Configuration{APIKey: honeybadgerApiKey})
		fmt.Println("enabled honeybadger integration")
	}

	honeybadger.Configure(honeybadger.Configuration{APIKey: "your api key"})

	caCertBytes, caCertBytesErr := getBytesFromEnv("CA_CERT")
	if caCertBytesErr != nil {
		errors = append(errors, caCertBytesErr)
	}

	serverCertBytes, serverCertBytesErr := getBytesFromEnv("SERVER_CERT")
	if serverCertBytesErr != nil {
		errors = append(errors, serverCertBytesErr)
	}

	serverKeyBytes, serverKeyBytesErr := getBytesFromEnv("SERVER_KEY")
	if serverKeyBytesErr != nil {
		errors = append(errors, serverKeyBytesErr)
	}

	smtpServer, smtpServerErr := getStringFromEnv("SMTP_SERVER")
	if smtpServerErr != nil {
		errors = append(errors, smtpServerErr)
	}

	smtpPort, smtpPortErr := getIntFromEnv("SMTP_PORT")
	if smtpPortErr != nil {
		errors = append(errors, smtpPortErr)
	}

	smtpUser, smtpUserErr := getStringFromEnv("SMTP_USER")
	if smtpUserErr != nil {
		errors = append(errors, smtpUserErr)
	}

	smtpPassword, smtpPasswordErr := getStringFromEnv("SMTP_PASSWORD")
	if smtpPasswordErr != nil {
		errors = append(errors, smtpPasswordErr)
	}

	fromDomain, fromDomainErr := getStringFromEnv("FROM_DOMAIN")
	if fromDomainErr != nil {
		errors = append(errors, fromDomainErr)
	}

	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	ConfigCaCertBytes = caCertBytes
	ConfigServerCertBytes = serverCertBytes
	ConfigServerKeyBytes = serverKeyBytes
	ConfigSmtpServer = smtpServer
	ConfigSmtpUser = smtpUser
	ConfigSmtpPassword = smtpPassword
	ConfigSmtpPort = smtpPort
	ConfigFromDomain = fromDomain
}

func loadTranslationsJson() {
	file, err := os.Open("translations.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(file)
	translations := make(map[string]interface{})
	err = decoder.Decode(&translations)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Translations = translations
}
