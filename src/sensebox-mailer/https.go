package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (mailer *senseBoxMailerServer) HandleSendRequest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var parsedRequest MailRequest
	err := decoder.Decode(&parsedRequest)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = SendMail(parsedRequest)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Fprint(w, "ok")
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
