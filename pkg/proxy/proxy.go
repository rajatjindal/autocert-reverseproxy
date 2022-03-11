package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//AutocertServer is autocert server
type AutocertServer struct {
	certCacheDir    string
	reverseProxyMap map[string]*httputil.ReverseProxy
}

func initReverseProxyMap(file string) (map[string]*httputil.ReverseProxy, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	reverseProxyMap := map[string]*httputil.ReverseProxy{}
	for host, upstream := range m {
		u, err := url.Parse(upstream)
		if err != nil {
			return nil, err
		}

		reverseProxyMap[host] = httputil.NewSingleHostReverseProxy(u)
	}

	return reverseProxyMap, nil
}

//Init inits the reverse proxy
func New(certDir, file string) (*AutocertServer, error) {
	reverseProxyMap, err := initReverseProxyMap(file)
	if err != nil {
		return nil, err
	}

	return &AutocertServer{
		certCacheDir:    certDir,
		reverseProxyMap: reverseProxyMap,
	}, nil
}

func (a *AutocertServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handling request for %s%s, %v\n", r.Host, r.URL.Path, r.Header)
	a.reverseProxyMap[r.Host].ServeHTTP(w, r)
}
