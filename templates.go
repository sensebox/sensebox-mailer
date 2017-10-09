package main

import (
	"bytes"
)

func prepareMailBody(templateName, language string, payload interface{}) (string, error) {
	var strBuffer bytes.Buffer
	err := templates.ExecuteTemplate(&strBuffer, templateName+"_"+language+".html", payload)
	if err != nil {
		return "", err
	}
	str := strBuffer.String()
	return str, nil
}
