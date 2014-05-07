package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"

	"stephensearles.com/muxchain"
	"stephensearles.com/muxchain/muxchainutil"
)

func main() {
	echoHandler := http.HandlerFunc(echo)
	authHandler := http.HandlerFunc(auth)

	muxchain.Chain("/", logMux(), muxchainutil.Gzip, echoHandler)
	muxchain.Chain("/noecho/", muxchainutil.Gzip, logMux())
	muxchain.Chain("/auth/", logMux(), muxchainutil.Gzip, authHandler, echoHandler)
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
	io.WriteString(w, req.URL.Path)
}

func auth(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("auth") == "yes" {
		w.Header().Add("X-Auth", "yes")
	} else {
		http.Error(w, "?auth=yes required", http.StatusForbidden)
	}
}

func Gzip(w http.ResponseWriter, req *http.Request, handlers ...http.Handler) {
	w.Header().Add("Content-Encoding", "gzip")
	w.Header().Add("Content-Type", "text/html")
	g, gw := NewGzipResponse(w)
	defer g.Close()
	muxchain.HandleChain(gw, req, handlers...)
}

type GzipResponse struct {
	http.ResponseWriter
	w *gzip.Writer
}

func NewGzipResponse(w http.ResponseWriter) (*gzip.Writer, http.ResponseWriter) {
	g := gzip.NewWriter(w)
	return g, &GzipResponse{w, g}
}

func (g *GzipResponse) Write(p []byte) (int, error) {
	return g.w.Write(p)
}

func (g *GzipResponse) Flush() error {
	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
	return g.w.Flush()
}
