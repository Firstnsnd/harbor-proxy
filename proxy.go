package proxy

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	Proxy  *httputil.ReverseProxy
	Target *url.URL
}

func NewReverseProxy(targetDomain string) (*ReverseProxy, error) {
	target, err := url.Parse(targetDomain)
	if err != nil {
		return nil, fmt.Errorf("error parsing target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	// Skip https cert
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // don't set true in prod environment
		},
	}
	return &ReverseProxy{Proxy: proxy, Target: target}, nil
}

func (rp *ReverseProxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = "http"
	}

	log.Println("Received request:", r.Method, r.URL.Path, r.Host, proto)
	host := r.Host

	//  ModifyResponse modify Header
	rp.Proxy.ModifyResponse = func(resp *http.Response) error {
		// Modify header path: v2/
		if resp.Header.Get("Www-Authenticate") != "" {
			// clear Header
			resp.Header.Del("Www-Authenticate")
			// setting new Header
			resp.Header.Set("Www-Authenticate", fmt.Sprintf(`Bearer realm="%s://%s/service/token",service="harbor-registry"`, proto, host))
		}

		// Modify  Location Header
		// for path: /uploads/   /manifests/
		if location := resp.Header.Get("Location"); location != "" {
			newLocation := strings.Replace(location, rp.Target.String(), fmt.Sprintf("%s://%s", proto, host), 1)
			resp.Header.Set("Location", newLocation)
		}
		return nil
	}
	r.URL.Scheme = rp.Target.Scheme
	r.URL.Host = rp.Target.Host
	r.Host = rp.Target.Host
	rp.Proxy.ServeHTTP(w, r)
}
