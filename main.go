package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var domain = "your_harbor_domain"

func main() {
	// target Harbor domain
	target, err := url.Parse(domain)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}
	// ReverseProxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Skip https cert
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // don't set true in prod environment
		},
	}

	// proxy handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http" // default use http
		}

		log.Println("Received request:", r.Method, r.URL.Path, r.Host, proto)
		host := r.Host
		//  ModifyResponse modify Header
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Modify header path: v2/
			log.Println("Received header:", resp.Header)

			if resp.Header.Get("Www-Authenticate") != "" {
				// clear Header
				resp.Header.Del("Www-Authenticate")
				// setting new Header

				resp.Header.Set("Www-Authenticate", fmt.Sprintf(`Bearer realm="%s://%s/service/token",service="harbor-registry"`, proto, host))
			}

			// Modify  Location Header
			// path: /uploads/   /manifests/
			if location := resp.Header.Get("Location"); location != "" {
				newLocation := strings.Replace(location, domain, fmt.Sprintf("%s://%s", proto, host), 1)
				resp.Header.Set("Location", newLocation)
			}
			log.Println("resp header:", resp.Header)

			return nil
		}

		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	})

	// start server
	log.Println("Starting proxy server on :8099")
	if err := http.ListenAndServe(":8099", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
