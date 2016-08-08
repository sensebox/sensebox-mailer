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
)

type MailerJSONResponse struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Request string `json:"request"`
}

func (mailer *senseBoxMailerServer) HandleSendRequest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var parsedRequests []MailRequest
	err := decoder.Decode(&parsedRequests)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// init data structure for response

	var jsonResponse []MailerJSONResponse

	hasError := false

	for _, request := range parsedRequests {
		currResponse := MailerJSONResponse{
			Status:  "ok",
			Request: request.Language + "_" + request.Template + "_" + request.Recipient.Address + "_" + time.Now().Format(time.RFC3339),
		}
		err = SendMail(request)
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
		fmt.Println("error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hasError == true {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)
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

	http.HandleFunc("/", mailer.HandleSendRequest)

	httpServer := &http.Server{
		Addr:      "0.0.0.0:3924",
		TLSConfig: tlsConfig,
	}
	fmt.Println("configured server")

	fmt.Println("starting server..")
	log.Println(httpServer.ListenAndServeTLS("", ""))
}
