package muxchain

import "net/http"

var DefaultMuxChain = &MuxChain{}

// Chain registers the pattern and http.Handler chain to the DefaultMuxChain.
func Chain(pattern string, handlers ...http.Handler) {
	DefaultMuxChain.Chain(pattern, handlers...)
}

type MuxChain struct {
	*http.ServeMux
}

// Chain registers a pattern to a sequence of http.Handlers. Upon receiving a request,
// the mux chain will find the best matching pattern and call that chain of handlers.
// The handlers will be called in turn until one of them writes a response or the end
// of the chain is reached.
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

func (m *MuxChain) handle(h http.Handler, lastHandler bool, w http.ResponseWriter, req *http.Request) bool {
	// Is the current handler a muxer?
	if childMux, ok := h.(muxer); ok {
		_, p := childMux.Handler(req)
		// Ignore this handler if this ServeMux doesn't apply, unless we have no more handlers
		if p == "" && !lastHandler {
			return true
		}
	}

	// Serve for the current handler

	cw := newChecked(w)
	h.ServeHTTP(cw.(http.ResponseWriter), req)
	return cw.IsWritten()
}

// muxer identifies types that act as a ServeMux.
type muxer interface {
	Handler(r *http.Request) (h http.Handler, pattern string)
}

type checkedResponseWriter interface {
	http.ResponseWriter
	IsWritten() bool
}

type checked struct {
	http.ResponseWriter
	written bool
}

func newChecked(w http.ResponseWriter) checkedResponseWriter {
	return &checked{w, false}
}

func (c *checked) Write(p []byte) (int, error) {
	c.written = true
	return c.ResponseWriter.Write(p)
}

func (c *checked) IsWritten() bool {
	return c.written
}
