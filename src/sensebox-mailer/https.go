package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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
		LogInfo("ServeHTTP", "Error:", err)
		http.Error(w, http.StatusText(status), status)
	}
}

func (mailer *senseBoxMailerServer) requestHandler(w http.ResponseWriter, req *http.Request) (int, error) {
	LogInfo("requestHandler", "incoming request")
	decoder := json.NewDecoder(req.Body)
	var parsedRequests []MailRequest
	err := decoder.Decode(&parsedRequests)
	if err != nil {
		LogInfo("requestHandler", "Error decoding JSON payload:", err)
		return http.StatusBadRequest, err
	}

	// init data structure for response
	var jsonResponse []MailerJSONResponse

	hasError := false
	requestTimestamp := time.Now().UTC().Unix()

	for _, request := range parsedRequests {
		request.Id = strconv.FormatInt(requestTimestamp, 36) + ";" + request.Language + ";" + request.Template + ";" + request.Recipient.Address
		currResponse := MailerJSONResponse{
			Status:  "ok",
			Request: request.Id,
		}
		err = mailer.SendMail(request)
		if err != nil {
			LogInfo("SendMail", "Error:", request.Id, err)
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
	LogInfo("StartHTTPSServer", "senseBox Mailer startup")

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(ConfigCaCertBytes); !ok {
		log.Fatalln("Unable to add CA certificate to client certificate pool")
		os.Exit(1)
	}
	LogInfo("StartHTTPSServer", "created client cert pool")

	myServerCertificate, err := tls.X509KeyPair(ConfigServerCertBytes, ConfigServerKeyBytes)
	if err != nil {
		log.Fatalln(err)
	}
	LogInfo("StartHTTPSServer", "imported server cert")

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		Certificates:             []tls.Certificate{myServerCertificate},
	}

	tlsConfig.BuildNameToCertificate()
	LogInfo("StartHTTPSServer", "built name to certificate")

	http.Handle("/", honeybadger.Handler(mailRequestHandler(mailer.requestHandler)))

	httpServer := &http.Server{
		Addr:      "0.0.0.0:3924",
		TLSConfig: tlsConfig,
	}
	LogInfo("StartHTTPSServer", "configured server")

	LogInfo("StartHTTPSServer", "starting server..")
	log.Fatal(httpServer.ListenAndServeTLS("", ""))
}
