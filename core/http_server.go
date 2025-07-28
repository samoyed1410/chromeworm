package core

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kgretzky/evilginx2/log"
)

type HttpServer struct {
	srv        *http.Server
	acmeTokens map[string]string
	tlsSrv     *http.Server
}

func NewHttpServer() (*HttpServer, error) {
	s := &HttpServer{}
	s.acmeTokens = make(map[string]string)

	r := mux.NewRouter()
	s.srv = &http.Server{
		Handler:      r,
		Addr:         ":80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.HandleFunc("/.well-known/acme-challenge/{token}", s.handleACMEChallenge).Methods("GET")
	r.PathPrefix("/").HandlerFunc(s.handleRedirect)

	return s, nil
}

func (s *HttpServer) Start() {
	go s.srv.ListenAndServe()
}

// Start HTTPS server with wildcard cert for all hostnames
func (s *HttpServer) StartTLS(config *Config) error {
	certPath := config.certificates.CertPath
	keyPath := config.certificates.KeyPath

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return &cert, nil
		},
		MinVersion: tls.VersionTLS12,
	}

	r := mux.NewRouter()
	r.HandleFunc("/.well-known/acme-challenge/{token}", s.handleACMEChallenge).Methods("GET")
	r.PathPrefix("/").HandlerFunc(s.handleRedirect)

	s.tlsSrv = &http.Server{
		Addr:         ":443",
		Handler:      r,
		TLSConfig:    tlsConfig,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := s.tlsSrv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Error("TLS server error: %v", err)
		}
	}()
	return nil
}

func (s *HttpServer) AddACMEToken(token string, keyAuth string) {
	s.acmeTokens[token] = keyAuth
}

func (s *HttpServer) ClearACMETokens() {
	s.acmeTokens = make(map[string]string)
}

func (s *HttpServer) handleACMEChallenge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	key, ok := s.acmeTokens[token]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Debug("http: found ACME verification token for URL: %s", r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(key))
}

func (s *HttpServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusFound)
}
