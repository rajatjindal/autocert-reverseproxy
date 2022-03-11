package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

func (a *AutocertServer) getHTTPServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)

	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":80",
		Handler:      mux,
	}

	return server
}

//Start starts the autocert server
func (a *AutocertServer) Start() {
	hostPolicy := func(ctx context.Context, requestedHost string) error {
		for host := range a.reverseProxyMap {
			if host == requestedHost {
				logrus.Infof("host policy matched for %s", requestedHost)
				return nil
			}
		}

		return fmt.Errorf("acme/autocert: %v host not allowed", requestedHost)
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(a.certCacheDir),
	}

	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":443",
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
		Handler:      m.HTTPHandler(a),
	}

	go func() {
		logrus.Infof("Starting HTTPS server on %s\n", server.Addr)
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			logrus.Fatalf("server.ListendAndServeTLS() failed with %s", err)
		}
	}()

	s := a.getHTTPServer()
	if m != nil {
		s.Handler = m.HTTPHandler(s.Handler)
	}

	logrus.Fatal(s.ListenAndServe())
}
