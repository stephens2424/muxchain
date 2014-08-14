package muxchainutil

import (
	"net/http"
	"strings"
)

// GlobMux muxes patterns with wildcard (*) components.
type GlobMux struct {
	paths map[string]http.Handler
}

// NewGlobMux initializes a new GlobMux.
func NewGlobMux() *GlobMux {
	return &GlobMux{paths: make(map[string]http.Handler)}
}

// Handler accepts a request and returns the appropriate handler for it, along
// with the pattern it matched. If no appropriate handler is found, the
// http.NotFoundHandler is returned along with an empty string. GlobMux will
// choose the handler that matches the most leading path components for the
// request.
func (g *GlobMux) Handler(req *http.Request) (h http.Handler, pattern string) {
	h, pattern = g.match(req.URL.Path)
	if pattern == "" {
		return http.NotFoundHandler(), pattern
	}
	return
}

// Handle registers a pattern to a handler.
func (g *GlobMux) Handle(pattern string, h http.Handler) {
	g.paths[pattern] = h
}

// ServeHTTP handles requests with the appropriate globbed handler.
func (g *GlobMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, _ := g.Handler(req)
	h.ServeHTTP(w, req)
}

func (g *GlobMux) match(path string) (h http.Handler, pattern string) {
	longestMatchLen := 0
	for matchPattern, handler := range g.paths {
		if pathMatch(matchPattern, path) {
			matchLen := strings.Count(matchPattern, "/")
			if matchLen > longestMatchLen {
				h = handler
				pattern = matchPattern
				longestMatchLen = matchLen
			}
		}
	}
	return
}

func pathMatch(pattern, path string) bool {
	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(pattern, "/")

	for i, patternPart := range patternParts {

		switch patternPart {
		case "*":
		default:
			if i > len(pathParts)-1 || pathParts[i] != patternPart {
				return false
			}
		}
	}
	return true
}
