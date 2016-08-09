package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/honeybadger-io/honeybadger-go"
)

type MailerJSONResponse struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Request string `json:"request"`
}

type mailRequestHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn mailRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		fmt.Println("error:", err)
		http.Error(w, http.StatusText(status), status)
	}
}

func (mailer *senseBoxMailerServer) requestHandler(w http.ResponseWriter, req *http.Request) (int, error) {
	decoder := json.NewDecoder(req.Body)
	var parsedRequests []MailRequest
	err := decoder.Decode(&parsedRequests)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// init data structure for response
	var jsonResponse []MailerJSONResponse

	hasError := false

	for _, request := range parsedRequests {
		currResponse := MailerJSONResponse{
			Status:  "ok",
			Request: request.Language + "_" + request.Template + "_" + request.Recipient.Address + "_" + time.Now().Format(time.RFC3339),
		}
		err = mailer.SendMail(request)
		if err != nil {
			fmt.Println("error:", err)
			currResponse.Status = "error"
			currResponse.Error = err.Error()
			hasError = true
		}
		jsonResponse = append(jsonResponse, currResponse)
	}

	jsonBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if hasError == true {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

	return http.StatusOK, nil
}

func (mailer *senseBoxMailerServer) StartHTTPSServer() {
	fmt.Println("senseBox Mailer startup")

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(ConfigCaCertBytes); !ok {
		log.Fatalln("Unable to add CA certificate to client certificate pool")
		os.Exit(1)
	}
	fmt.Println("created client cert pool")

	myServerCertificate, err := tls.X509KeyPair(ConfigServerCertBytes, ConfigServerKeyBytes)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("imported server cert")

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		Certificates:             []tls.Certificate{myServerCertificate},
	}

	tlsConfig.BuildNameToCertificate()
	fmt.Println("built name to certificate")

	http.Handle("/", honeybadger.Handler(mailRequestHandler(mailer.requestHandler)))

	httpServer := &http.Server{
		Addr:      "0.0.0.0:3924",
		TLSConfig: tlsConfig,
	}
	fmt.Println("configured server")

	fmt.Println("starting server..")
	log.Fatal(httpServer.ListenAndServeTLS("", ""))
}
