package mailer

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(status), err.Error()), status)
	}
}

func (mailer *MailerServer) requestHandler(w http.ResponseWriter, req *http.Request) (int, error) {
	LogInfo("requestHandler", "incoming request")

	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		LogInfo("requestHandler", "Error reading request body:", err)
		return http.StatusBadRequest, fmt.Errorf("Error reading request body: %w", err)
	}
	defer req.Body.Close()

	var parsedRequests []MailRequest
	err = json.Unmarshal(reqBytes, &parsedRequests)
	if err != nil {
		LogInfo("requestHandler", "Error decoding JSON payload:", err)
		return http.StatusBadRequest, fmt.Errorf("Error decoding JSON payload: %w", err)
	}

	// init data structure for response
	var jsonResponse []MailerJSONResponse

	httpCode := http.StatusOK

	for _, request := range parsedRequests {
		currResponse := MailerJSONResponse{
			Status:  "ok",
			Request: request.ID,
		}
		err = mailer.sendMail(request)
		if err != nil {
			LogInfo("SendMail", "Error:", request.ID, err)
			currResponse.Status = "error"
			currResponse.Error = err.Error()
			httpCode = http.StatusBadRequest
		}
		jsonResponse = append(jsonResponse, currResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	return httpCode, json.NewEncoder(w).Encode(jsonResponse)
}

func (mailer *MailerServer) startHTTPSServer() error {
	LogInfo("StartHTTPSServer", "senseBox Mailer startup")

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(mailer.CaCert); !ok {
		return fmt.Errorf("Unable to add CA certificate to client certificate pool")
	}
	LogInfo("StartHTTPSServer", "created client cert pool")

	myServerCertificate, err := tls.X509KeyPair(mailer.ServerCert, mailer.ServerKey)
	if err != nil {
		return err
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

	http.Handle("/", mailRequestHandler(mailer.requestHandler))

	httpServer := &http.Server{
		Addr:      "0.0.0.0:3924",
		TLSConfig: tlsConfig,
	}
	LogInfo("StartHTTPSServer", "configured server")

	LogInfo("StartHTTPSServer", "starting server on address 0.0.0.0:3924")
	log.Fatal(httpServer.ListenAndServeTLS("", ""))

	return nil
}
