package main

import (
	"fmt"
	"os"
	"strconv"
)

const envPrefix = "SENSEBOX_MAILER_"

func initConfigFromEnv() (caCert, serverCert, serverKey []byte, smtpServer, smtpUser, smtpPassword, fromDomain string, smtpPort int, errors []error) {
	errors = make([]error, 0)

	caCert, caCertBytesErr := getBytesFromEnv("CA_CERT")
	if caCertBytesErr != nil {
		errors = append(errors, caCertBytesErr)
	}

	serverCert, serverCertBytesErr := getBytesFromEnv("SERVER_CERT")
	if serverCertBytesErr != nil {
		errors = append(errors, serverCertBytesErr)
	}

	serverKey, serverKeyBytesErr := getBytesFromEnv("SERVER_KEY")
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
		return
	}
	return
}

func getStringFromEnv(key string) (string, error) {
	str := os.Getenv(envPrefix + key)
	if len(str) == 0 {
		return "", fmt.Errorf("Please add %s%s to your environment", envPrefix, key)
	}
	return str, nil
}

func getBytesFromEnv(key string) ([]byte, error) {
	str, err := getStringFromEnv(key)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func getIntFromEnv(key string) (int, error) {
	str, err := getStringFromEnv(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("Environment key %s%s is not parseable as integer", envPrefix, key)
	}
	return i, nil
}
