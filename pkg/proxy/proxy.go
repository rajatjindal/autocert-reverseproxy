package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

//AutocertServer is autocert server
type AutocertServer struct {
	HTTPPort     string
	HTTPSPort    string
	AllowedHost  map[string]bool
	ReverseProxy *httputil.ReverseProxy

	sync.Mutex
}

//New returns new autocert server
func New() (*AutocertServer, error) {
	u, _ := url.Parse("http://localhost:8083")
	as := &AutocertServer{
		HTTPPort:  ":80",
		HTTPSPort: ":443",
		AllowedHost: map[string]bool{
			"gapiv4.rajatjindal.com": true,
			"gapiv5.rajatjindal.com": true,
			"gapiv6.rajatjindal.com": true,
		},
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
	}

	return as, nil
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
		Addr:         a.HTTPPort,
		Handler:      mux,
	}

	return server
}

//Start starts the autocert server
func (a *AutocertServer) Start() {
	hostPolicy := func(ctx context.Context, host string) error {
		for k, v := range a.AllowedHost {
			if v && k == host {
				logrus.Infof("host policy matched for %s", host)
				return nil
			}
		}

		return fmt.Errorf("acme/autocert: only %v host is allowed", a.AllowedHost)
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache("/home/app/cert"),
	}

	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         a.HTTPSPort,
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
		Handler:      m.HTTPHandler(a.ReverseProxy),
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
