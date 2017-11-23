package main

import (
	"bytes"
)

func prepareMailBody(templateName, language string, payload interface{}) ([]byte, error) {
	var strBuffer bytes.Buffer
	err := templates.ExecuteTemplate(&strBuffer, templateName+"_"+language+".html", payload)
	if err != nil {
		return []byte{}, err
	}
	str := strBuffer.Bytes()
	return str, nil
}
