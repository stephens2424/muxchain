package muxchainutil

import (
	"fmt"
	"net/http"
	"strings"
)

// MethodMux allows the caller to specify handlers that are specific to a particular HTTP method
type MethodMux struct {
	methods         map[string]*http.ServeMux
	NotFoundHandler http.Handler
}

// NewMethodMux returns a new MethodMux
func NewMethodMux() *MethodMux {
	methods := make(map[string]*http.ServeMux)
	methods["*"] = http.NewServeMux()
	return &MethodMux{methods, http.HandlerFunc(NoopHandler)}
}

func NoopHandler(w http.ResponseWriter, req *http.Request) {
}

// Handle registers a pattern to a particular handler. The pattern may optionally begin with an HTTP
// method, followed by a space, e.g.: "GET /homepage". The method may be *, which matches all methods.
// The method may also be omitted, which is the same as *.
func (m *MethodMux) Handle(pattern string, h http.Handler) {
	s := strings.SplitN(pattern, " ", 2)
	if len(s) < 2 {
		m.handle("*", s[0], h)
	} else {
		m.handle(s[0], s[1], h)
	}
}

func (m *MethodMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, _ := m.Handler(req)
	h.ServeHTTP(w, req)
}

// HandleMethods registers a pattern to a handler for the given methods.
func (m *MethodMux) HandleMethods(pattern string, h http.Handler, methods ...string) {
	for _, method := range methods {
		m.Handle(fmt.Sprintf("%s %s", method, pattern), h)
	}
}

func (m *MethodMux) handle(method string, pattern string, h http.Handler) {
	muxer, ok := m.methods[method]
	if !ok {
		muxer = http.NewServeMux()
		m.methods[method] = muxer
	}
	muxer.Handle(pattern, h)
}

// Handler selects a handler for a request. It will first attempt to chose a handler
// that matches the particular method. If a handler is not found, pattern will be
// the empty string.
func (m *MethodMux) Handler(req *http.Request) (h http.Handler, pattern string) {
	h, pattern = m.handleMethod(req.Method, req)
	if pattern != "" {
		return
	}
	return m.handleMethod(req.Method, req)
}

// handleMethod selects a handler for a request for a particular method, ignoring the method
// actually on the request.
func (m *MethodMux) handleMethod(method string, req *http.Request) (h http.Handler, pattern string) {
	muxer, ok := m.methods[method]
	if !ok {
		return m.NotFoundHandler, ""
	}
	return muxer.Handler(req)
}
