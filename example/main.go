package main

import (
	"github.com/Firstnsnd/harbor-proxy" // 替换为实际的模块路径
	"log"
	"net/http"
)

var domain = "your_harbor_domain"

func main() {
	reverseProxy, err := proxy.NewReverseProxy(domain)
	if err != nil {
		log.Fatalf("Error creating reverse proxy: %v", err)
	}

	http.HandleFunc("/", reverseProxy.HandleRequest)

	log.Println("Starting proxy server on :8099")
	if err := http.ListenAndServe(":8099", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
