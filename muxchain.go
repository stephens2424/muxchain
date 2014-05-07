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

func (m *MuxChain) ServeHTTPChain(w http.ResponseWriter, req *http.Request, handlers ...http.Handler) {
	HandleChain(w, req, handlers...)
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
		HandleChain(w, req, handlers...)
	})
}

// HandleChain is the utility function chained handlers are responsible for calling
// when they are complete.
func HandleChain(w http.ResponseWriter, req *http.Request, handlers ...http.Handler) {
	for i, h := range handlers {
		var remaining []http.Handler
		if len(handlers) > i+1 {
			remaining = handlers[i+1:]
		}
		if handle(h, remaining, w, req) {
			return
		}
	}
}

// ChainedHander allows implementers to call the handlers after them in a muxchain
// on their own. This allows handlers to defer functions until after the handler
// chain following them has been completed.
type ChainedHandler interface {
	http.Handler

	// ServeHTTPChain allows the handler to do its work and then call the remaining
	// handler chain. Implementers should call muxchain.HandleChain when complete.
	ServeHTTPChain(w http.ResponseWriter, req *http.Request, h ...http.Handler)
}

// ChainedHandlerFunc represents a handler that is able to be chained to subsequent handlers.
type ChainedHandlerFunc func(w http.ResponseWriter, req *http.Request, handlers ...http.Handler)

// ServeHTTP allows ChainedHandlerFuncs to implement http.Handler
func (c ChainedHandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c(w, req, nil)
}

// ServeHTTPChain allows ChainedHandlerFuncs to be identified
func (c ChainedHandlerFunc) ServeHTTPChain(w http.ResponseWriter, req *http.Request, handlers ...http.Handler) {
	c(w, req, handlers...)
}

// handle runs and attempts to serve on the current handler. It returns true if data
// was written to the response writer.
func handle(h http.Handler, remaining []http.Handler, w http.ResponseWriter, req *http.Request) bool {
	// Is the current handler a Muxer?
	if childMux, ok := h.(Muxer); ok {
		_, p := childMux.Handler(req)
		// Ignore this handler if this ServeMux doesn't apply, unless we have no more handlers
		if p == "" && len(remaining) == 0 {
			return true
		}
	}

	cw := &checked{w, false}
	if chainedHandler, ok := h.(ChainedHandler); ok {
		chainedHandler.ServeHTTPChain(cw, req, remaining...)
		return true
	}

	// Serve for the current handler
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
