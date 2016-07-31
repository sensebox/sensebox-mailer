package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

const ENV_PREFIX = "SENSEBOX_MAILER_"

type senseBoxMailer struct{}

// HelloUser is a view that greets a user
func (mailer *senseBoxMailer) HelloUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello %v! \n", req.TLS.PeerCertificates[0].Subject)
}

func (mailer *senseBoxMailer) StartServer() {
	fmt.Println("senseBox Mailer startup")

	caCertBytes := getBytesFromEnvOrFail("CA_CERT")
	serverCertBytes := getBytesFromEnvOrFail("SERVER_CERT")
	serverKeyBytes := getBytesFromEnvOrFail("SERVER_KEY")

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(caCertBytes); !ok {
		log.Fatalln("Unable to add CA certificate to client certificate pool")
		os.Exit(1)
	}
	fmt.Println("created client cert pool")

	myServerCertificate, err := tls.X509KeyPair(serverCertBytes, serverKeyBytes)
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

	http.HandleFunc("/", mailer.HelloUser)

	httpServer := &http.Server{
		Addr:      "0.0.0.0:3924",
		TLSConfig: tlsConfig,
	}
	fmt.Println("configured server")

	fmt.Println("starting server..")
	log.Println(httpServer.ListenAndServeTLS("", ""))
}

func main() {

	mailer := senseBoxMailer{}
	mailer.StartServer()
}
