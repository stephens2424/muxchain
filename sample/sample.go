package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"stephensearles.com/muxchain"
	"stephensearles.com/muxchain/muxchainutil"
)

func main() {
	echoHandler := http.HandlerFunc(echo)
	authHandler := http.HandlerFunc(auth)

	pathHandler := muxchainutil.NewPathMuxer()
	pathHandler.Handle("/id/:id", echoHandler)

	muxchain.Chain("/", logMux(), muxchainutil.Gzip, echoHandler)
	muxchain.Chain("/noecho/", muxchainutil.Gzip, logMux())
	muxchain.Chain("/auth/", logMux(), muxchainutil.Gzip, authHandler, echoHandler)
	muxchain.Chain("/id/", logMux(), muxchainutil.Gzip, pathHandler)
	http.ListenAndServe(":36363", muxchain.Default)
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
	w.Header().Add("Content-Type", "text/html")
	if id := req.FormValue("id"); id != "" {
		fmt.Fprintf(w, "Your id is %s", id)
		return
	}
	io.WriteString(w, req.URL.Path)
}

func auth(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("auth") == "yes" {
		w.Header().Add("X-Auth", "yes")
	} else {
		http.Error(w, "?auth=yes required", http.StatusForbidden)
	}
}
