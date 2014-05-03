package muxchain

import "net/http"

var Default = &MuxChain{}

// Chain registers the pattern and http.Handler chain to the DefaultMuxChain.
func Chain(pattern string, handlers ...http.Handler) {
	Default.Chain(pattern, handlers...)
}

type MuxChain struct {
	*http.ServeMux
}

// Chain registers a pattern to a sequence of http.Handlers. Upon receiving a request,
// the mux chain will find the best matching pattern and call that chain of handlers.
// The handlers will be called in turn until one of them writes a response or the end
// of the chain is reached. If one of the handlers is a Muxer (e.g. http.ServeMux),
// the MuxChain will skip it if no pattern matches.
func (m *MuxChain) Chain(pattern string, handlers ...http.Handler) {
	if m.ServeMux == nil {
		m.ServeMux = http.NewServeMux()
	}
	m.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if len(handlers) == 0 {
			http.NotFound(w, req)
		}
		for i, h := range handlers {
			if m.handle(h, i < len(handlers)-1, w, req) {
				return
			}
		}
	})
}

// handle runs and attempts to serve on the current handler. It returns true if data
// was written to the response writer.
func (m *MuxChain) handle(h http.Handler, lastHandler bool, w http.ResponseWriter, req *http.Request) bool {
	// Is the current handler a Muxer?
	if childMux, ok := h.(Muxer); ok {
		_, p := childMux.Handler(req)
		// Ignore this handler if this ServeMux doesn't apply, unless we have no more handlers
		if p == "" && !lastHandler {
			return true
		}
	}

	// Serve for the current handler
	cw := &checked{w, false}
	h.ServeHTTP(cw, req)
	return cw.written
}

// Muxer identifies types that act as a ServeMux.
type Muxer interface {
	Handler(r *http.Request) (h http.Handler, pattern string)
}

type checked struct {
	http.ResponseWriter
	written bool
}

func (c *checked) Write(p []byte) (int, error) {
	c.written = true
	return c.ResponseWriter.Write(p)
}
