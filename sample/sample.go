package main

import (
	"io"
	"log"
	"net/http"

	"stephensearles.com/muxchain"
)

func main() {
	echoHandler := http.HandlerFunc(echo)
	authHandler := http.HandlerFunc(auth)

	muxchain.Chain("/", logMux(), echoHandler)
	muxchain.Chain("/noecho/", logMux())
	muxchain.Chain("/auth/", logMux(), authHandler, echoHandler)
	http.ListenAndServe(":36363", muxchain.DefaultMuxChain)
}

func logMux() *http.ServeMux {
	l := http.NewServeMux()
	l.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL.Path)
	})
	l.HandleFunc("/private", func(w http.ResponseWriter, req *http.Request) {
		// do nothing
	})
	return l
}

func echo(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, req.URL.Path)
}

func auth(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("auth") != "yes" {
		http.Error(w, "?auth=yes required", http.StatusForbidden)
	}
}
