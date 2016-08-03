package main

import (
	"html/template"
)

type senseBoxMailerServer struct {
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	initConfigFromEnv()
	loadTranslationsJson()
	templates.Option("missingkey=error")

	mailer := senseBoxMailerServer{}
	mailer.StartHTTPSServer()
}
