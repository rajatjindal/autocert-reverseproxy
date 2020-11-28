package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

//AutocertServer is autocert server
type AutocertServer struct {
	HTTPAddr     string
	HTTPSAddr    string
	AllowedHost  []string
	CertCacheDir string
	Upstream     string

	reverseProxy *httputil.ReverseProxy
}

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
		Addr:         a.HTTPAddr,
		Handler:      mux,
	}

	return server
}

//Init inits the reverse proxy
func (a *AutocertServer) Init() error {
	if len(a.AllowedHost) == 0 {
		return fmt.Errorf("allowed-host list is empty. need atleast one entry")
	}

	u, err := url.Parse(a.Upstream)
	if err != nil {
		return err
	}

	a.reverseProxy = httputil.NewSingleHostReverseProxy(u)
	return nil
}

//Start starts the autocert server
func (a *AutocertServer) Start() {
	hostPolicy := func(ctx context.Context, host string) error {
		for _, v := range a.AllowedHost {
			if v == host {
				logrus.Debugf("host policy matched for %s", host)
				return nil
			}
		}

		return fmt.Errorf("acme/autocert: only %v host is allowed", a.AllowedHost)
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(a.CertCacheDir),
	}

	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         a.HTTPSAddr,
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
		Handler:      m.HTTPHandler(a.reverseProxy),
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
