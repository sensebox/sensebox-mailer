package mailer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"time"

	"github.com/sensebox/sensebox-mailer/mailer/templates"

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
// MailerServer wraps the https server and the SMTP daemon
type MailerServer struct {
	daemon                                         chan *email.Email
	CaCert, ServerCert, ServerKey                  []byte
	SMTPServer, SMTPUser, SMTPPassword, FromDomain string
	SMTPPort                                       int
}

type MailRequestAttachment struct {
	Filename       string `json:"filename"`
	DecodedContent []byte `json:"contents"`
}

func (a *MailRequestAttachment) UnmarshalJSON(jsonBytes []byte) error {
	var attachment map[string]string

	if err := json.Unmarshal(jsonBytes, &attachment); err != nil {
		return err
	}
	if attachment["filename"] == "" {
		return errors.New("key 'filename' of key 'attachment' is required")
	}
	if attachment["contents"] == "" {
		return errors.New("key 'contents' of key 'attachment' is required")
	}

	data, err := base64.StdEncoding.DecodeString(attachment["contents"])
	if err != nil {
		return err
	}

	*a = MailRequestAttachment{
		Filename:       attachment["filename"],
		DecodedContent: data,
	}

	return nil
}

type MailRequestEmailAddress struct {
	Address       string `json:"address"`
	Name          string `json:"name"`
	ParsedAddress *mail.Address
}

func (a *MailRequestEmailAddress) UnmarshalJSON(jsonBytes []byte) error {
	var ma map[string]string

	if err := json.Unmarshal(jsonBytes, &ma); err != nil {
		return err
	}
	if ma["address"] == "" {
		return errors.New("key 'address' of key 'recipient' is required")
	}
	if ma["name"] == "" {
		return errors.New("key 'name' of key 'recipient' is required")
	}
	addr, err := mail.ParseAddress(fmt.Sprintf("%s <%s>", ma["name"], ma["address"]))
	if err != nil {
		return err
	}

	*a = MailRequestEmailAddress{
		Address:       ma["address"],
		Name:          ma["name"],
		ParsedAddress: addr,
	}

	return nil
}

type MailRequest struct {
	Recipient  MailRequestEmailAddress `json:"recipient"`
	Payload    map[string]interface{}  `json:"payload"`
	Attachment *MailRequestAttachment  `json:"attachment,omitempty"`
	Body       []byte
	FromName   string
	Subject    string
	ID         string
}

func (mr *MailRequest) UnmarshalJSON(jsonBytes []byte) error {
	var request map[string]*json.RawMessage
	var templateName, language string

	if err := json.Unmarshal(jsonBytes, &request); err != nil {
		return err
	}
	// check if the required keys are there
	if request["lang"] == nil {
		return errors.New("key 'lang' is required")
	}
	if request["template"] == nil {
		return errors.New("key 'template' is required")
	}

	if err := json.Unmarshal(*request["template"], &templateName); err != nil {
		return err
	}

	if err := json.Unmarshal(*request["lang"], &language); err != nil {
		return err
	}

	// check if the requested template is available
	templ, err := templates.GetTemplate(templateName, language)
	if err != nil {
		return err
	}

	// execute the template
	var payload map[string]interface{}
	if err := json.Unmarshal(*request["payload"], &payload); err != nil {
		return err
	}
	templateBytes, err := templ.Execute(payload)
	if err != nil {
		return err
	}

	// recipient
	var recipient MailRequestEmailAddress
	if err := json.Unmarshal(*request["recipient"], &recipient); err != nil {
		return err
	}

	*mr = MailRequest{
		ID: fmt.Sprintf("%s;%s;%s;%s",
			strconv.FormatInt(time.Now().UTC().UnixNano(), 36),
			language, templateName, recipient.Address),
		Recipient: recipient,
		Body:      templateBytes,
		FromName:  templ.FromName,
		Subject:   templ.Subject,
	}

	// attachment is optional
	if request["attachment"] != nil {
		var attachment *MailRequestAttachment
		if err := json.Unmarshal(*request["attachment"], &attachment); err != nil {
			return err
		}
		mr.Attachment = attachment

	}
	return nil
}

func (mailer *MailerServer) Start() error {
	ctTemplates, err := templates.FromJSON()
	if err != nil {
		return err
	}
	LogInfo("Start MailerServer", "Imported", ctTemplates, "templates")

	mailer.startMailerDaemon()
	defer close(mailer.daemon)

	err = mailer.startHTTPSServer()
	if err != nil {
		return err
	}

	return nil
}

func (mailer *MailerServer) sendMail(req MailRequest) error {
	headers := textproto.MIMEHeader{}
	headers.Add("senseBoxMailerInternalId", req.ID)

	m := &email.Email{
		To:      []string{req.Recipient.ParsedAddress.String()},
		From:    fmt.Sprintf("%s <support@%s>", req.FromName, mailer.FromDomain),
		Subject: req.Subject,
		HTML:    req.Body,
		Headers: headers,
	}

	if req.Attachment != nil {
		_, err := m.Attach(bytes.NewReader(req.Attachment.DecodedContent), req.Attachment.Filename, "text/plain")
		if err != nil {
			return err
		}
	}
	LogInfo("SendMail", req.ID, "submitting mail to mailer daemon")
	mailer.daemon <- m

	return nil
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func (mailer *MailerServer) startMailerDaemon() {
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
						fmt.Sprintf("%s:%d", mailer.SMTPServer, mailer.SMTPPort),
						4,
						unencryptedAuth{smtp.PlainAuth("", mailer.SMTPUser, mailer.SMTPPassword, mailer.SMTPServer)},
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

	mailer.daemon = ch
}
