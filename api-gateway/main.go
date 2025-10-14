package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	http.HandleFunc("/api/", handleProxy)

	log.Println("API Gateway escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	var target *url.URL

	// Ruteo según prefijo
	if strings.HasPrefix(r.URL.Path, "/api/users") {
		target, _ = url.Parse("http://localhost:8081")
	} else if strings.HasPrefix(r.URL.Path, "/api/activities") {
		target, _ = url.Parse("http://localhost:8082")
	} else if strings.HasPrefix(r.URL.Path, "/api/search") {
		target, _ = url.Parse("http://localhost:8083")
	} else {
		http.NotFound(w, r)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Quitar el prefijo /api al reenviar
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
	proxy.ServeHTTP(w, r)
}
