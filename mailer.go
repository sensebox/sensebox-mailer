package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	// "github.com/jordan-wright/email"
	"github.com/lovego/email"
)

/*
{
  "template": "registration",
  "lang": "en",
  "recipient": {
	"address": "email@address.com",
	"name": "Philip J. Fry"
  },
  "payload": {
    "user": {
      "firstname": "Philip J.",
      "lastname": "Fry",
      "apikey": "<some valid apikey>"
    },
    "box": {
      "id": "<some valid senseBox id>",
      "sensors": [
        {
          "title": "<some title>",
          "type": "<some type>",
          "id": "<some valid senseBox sensor id>"
        },
        ...
      ]
    }
  },
  "attachment": {
    "filename": "senseBox.ino",
    "contents": "<file contents in base64>"
  }
}
*/

var addressParser = mail.AddressParser{}

type MailRequestAttachment struct {
	Filename string `json:"filename"`
	Contents string `json:"contents"`
}

type MailRequestDecodedAttachment struct {
	Filename string
	Contents []byte
}

type MailRequestEmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type MailRequest struct {
	Template          string                  `json:"template"`
	Language          string                  `json:"lang"`
	Recipient         MailRequestEmailAddress `json:"recipient"`
	Payload           map[string]interface{}  `json:"payload"`
	Attachment        *MailRequestAttachment  `json:"attachment,omitempty"`
	DecodedAttachment *MailRequestDecodedAttachment
	BuiltTemplate     []byte
	EmailFrom         MailRequestEmailAddress
	Subject           string
	Id                string
}

func (request *MailRequest) validateAndParseRequest() error {
	// check if the required keys are there
	if request.Language == "" {
		return errors.New("key 'lang' is required")
	}
	if request.Template == "" {
		return errors.New("key 'template' is required")
	}
	if request.Recipient.Address == "" {
		return errors.New("key 'address' of key 'recipient' is required")
	}
	if request.Recipient.Name == "" {
		return errors.New("key 'name' of key 'recipient' is required")
	}

	// check if the supplied language is just an alias for another language
	for avaliableLang := range Translations {
		if strings.HasPrefix(request.Language, avaliableLang) {
			LogInfo("validateAndParseRequest", request.Id, "Converting "+request.Language+" to "+avaliableLang)
			request.Language = avaliableLang
			break
		}
	}

	// check if the supplied language is available
	_, present := Translations[request.Language]
	if present == false {
		LogInfo("validateAndParseRequest", request.Id, "Language "+request.Language+" not found. Falling back to 'en'")
		request.Language = "en"
	}

	// parse the Recipients address
	_, err := addressParser.Parse(request.Recipient.Address)
	if err != nil {
		return err
	}

	// decode the Attachment
	if request.Attachment != nil {
		data, err := base64.StdEncoding.DecodeString(request.Attachment.Contents)
		if err != nil {
			return err
		}
		decodedAttachment := MailRequestDecodedAttachment{
			Filename: request.Attachment.Filename,
			Contents: data,
		}
		request.DecodedAttachment = &decodedAttachment
	}

	// Fill FromAddress and SenderName
	senderName, err := getTranslation(request.Language, request.Template, "fromName")
	if err != nil {
		return err
	}
	request.EmailFrom = MailRequestEmailAddress{
		Address: "no-reply@" + ConfigFromDomain,
		Name:    senderName,
	}

	// Fill in Subject
	subj, err := getTranslation(request.Language, request.Template, "subject")
	if err != nil {
		return err
	}
	request.Subject = subj

	// execute the template
	s, err := prepareMailBody(request.Template, request.Language, request.Payload)
	if err != nil {
		return err
	}
	request.BuiltTemplate = s

	return nil
}

func (mailer *senseBoxMailerServer) SendMail(req MailRequest) error {

	err := req.validateAndParseRequest()
	if err != nil {
		return err
	}

	headers := textproto.MIMEHeader{}
	headers.Add("senseBoxMailerInternalId", req.Id)

	m := &email.Email{
		To:      []string{fmt.Sprintf("%s <%s>", req.Recipient.Name, req.Recipient.Address)},
		From:    fmt.Sprintf("%s <%s>", req.EmailFrom.Name, req.EmailFrom.Address),
		Subject: req.Subject,
		HTML:    req.BuiltTemplate,
		Headers: headers,
	}

	if req.DecodedAttachment != nil {
		_, err := m.Attach(bytes.NewReader(req.DecodedAttachment.Contents), req.DecodedAttachment.Filename, "text/plain")
		if err != nil {
			return err
		}
	}
	LogInfo("SendMail", req.Id, "submitting mail to mailer daemon")
	mailer.Daemon <- m

	return nil
}

func (mailer *senseBoxMailerServer) startMailerDaemon() {
	ch := make(chan *email.Email)

	go func() {

		var d *email.Pool
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				if !open {
					LogInfo("mailerDaemon", "trying to open connection to SMTP server")
					d = email.NewPool(
						fmt.Sprintf("%s:%d", ConfigSmtpServer, ConfigSmtpPort),
						4,
						smtp.PlainAuth("", ConfigSmtpUser, ConfigSmtpPassword, ConfigSmtpServer),
					)
					open = true
					LogInfo("mailerDaemon", "successfully opened connection to SMTP server")
				}
				LogInfo("mailerDaemon", m.Headers.Get("senseBoxMailerInternalId"), "trying to send mail")
				if err := d.Send(m, 5*time.Second); err != nil {
					LogInfo("mailerDaemon", "Error for", m.Headers.Get("senseBoxMailerInternalId"), err)
				}
				LogInfo("mailerDaemon", m.Headers.Get("senseBoxMailerInternalId"), "mail submitted to SMTP server")
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					LogInfo("mailerDaemon", "trying to close connection to SMTP server")
					d.Close()
					open = false
					LogInfo("mailerDaemon", "closed connection to SMTP server after 30 seconds of inactivity")
				}
			}
		}
	}()

	mailer.Daemon = ch
}
