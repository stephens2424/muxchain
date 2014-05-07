package muxchainutil

import (
	"compress/gzip"
	"net/http"

	"stephensearles.com/muxchain"
)

var Gzip = muxchain.ChainedHandlerFunc(gzipHandler)

func gzipHandler(w http.ResponseWriter, req *http.Request, h ...http.Handler) {
	w.Header().Add("Content-Encoding", "gzip")
	w.Header().Add("Content-Type", "text/html")
	g, gw := newGzipResponse(w)
	defer g.Close()
	muxchain.HandleChain(gw, req, h...)
}

type gzipResponse struct {
	http.ResponseWriter
	w *gzip.Writer
}

func newGzipResponse(w http.ResponseWriter) (*gzip.Writer, http.ResponseWriter) {
	g := gzip.NewWriter(w)
	return g, &gzipResponse{w, g}
}

func (g *gzipResponse) Write(p []byte) (int, error) {
	return g.w.Write(p)
}

func (g *gzipResponse) Flush() error {
	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
	return g.w.Flush()
}
