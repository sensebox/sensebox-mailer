package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
	_ "github.com/sensebox/sensebox-mailer/statik"
)

var statikFS http.FileSystem
var theTemplates []LocalizedTemplate

type LocalizedTemplate struct {
	Language     string `json:"language"`
	TemplateName string `json:"template"`
	FromName     string `json:"fromName"`
	Subject      string `json:"subject"`
	Template     *template.Template
}

func (td *LocalizedTemplate) UnmarshalJSON(jsonBytes []byte) error {
	var t map[string]string

	if err := json.Unmarshal(jsonBytes, &t); err != nil {
		return err
	}
	if err := initStatikFS(); err != nil {
		return err
	}

	f, err := statikFS.Open(fmt.Sprintf("/%s_%s.html", t["template"], t["language"]))
	if err != nil {
		return err
	}
	templateBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	template, err := template.New(t["template"]).Parse(string(templateBytes))
	if err != nil {
		return err
	}
	template.Option("missingkey=error")

	*td = LocalizedTemplate{
		t["language"],
		t["template"],
		t["fromName"],
		t["subject"],
		template,
	}

	return nil
}

func initStatikFS() error {
	sFS, err := fs.New()
	if err != nil {
		return err
	}
	statikFS = sFS
	return nil
}

// FromJSON initializes the templates from the templates/templates.json file
func FromJSON() (int, error) {
	if err := initStatikFS(); err != nil {
		return 0, err
	}

	f, err := statikFS.Open("/templates.json")
	if err != nil {
		return 0, err
	}
	jsonBytes, err := ioutil.ReadAll(f)

	if err != nil {
		return 0, err
	}
	f.Close()

	err = json.Unmarshal(jsonBytes, &theTemplates)
	if err != nil {
		return 0, err
	}

	return len(theTemplates), nil
}

// GetTemplate returns the template matching the templateName and the language
func GetTemplate(templateName, language string) (LocalizedTemplate, error) {
	if theTemplates == nil {
		_, err := FromJSON()
		if err != nil {
			return LocalizedTemplate{}, err
		}
	}

	// check for direct or prefix (de_DE -> de) match
	for _, template := range theTemplates {
		if template.TemplateName == templateName && (template.Language == language || strings.HasPrefix(language, template.Language) == true) {
			return template, nil
		}
	}

	return LocalizedTemplate{}, fmt.Errorf("Template not available (template: '%s', language: '%s')", templateName, language)
}

// Execute executes the Template, filling it with the values from the payload
func (td LocalizedTemplate) Execute(payload interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	err := td.Template.Execute(&buffer, payload)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}
