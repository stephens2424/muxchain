package muxchainutil

import (
	"net/http"
	"strings"
)

type PathMuxer struct {
	*GlobMuxer
	patternVariables map[string]map[int]string // maps pattern component index to variable name
}

func NewPathMuxer() *PathMuxer {
	return &PathMuxer{NewGlobMuxer(), make(map[string]map[int]string)}
}

func (p *PathMuxer) Handle(pattern string, h http.Handler) {
	patternParts := strings.Split(pattern, "/")
	variables := make(map[int]string)
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			variables[i] = part[1:]
			patternParts[i] = "*"
		}
	}
	pattern = strings.Join(patternParts, "/")
	p.patternVariables[pattern] = variables
	p.GlobMuxer.Handle(pattern, h)
}

func (p *PathMuxer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	h, pattern := p.Handler(req)
	if variables, hasVariables := p.patternVariables[pattern]; hasVariables {
		for i, part := range strings.Split(req.URL.Path, "/") {
			if varname, ok := variables[i]; ok {
				req.URL.Query().Set(varname, part)
				req.Form.Set(varname, part)
			}
		}
	}
	h.ServeHTTP(w, req)
}

type GlobMuxer struct {
	paths map[string]http.Handler
}

func NewGlobMuxer() *GlobMuxer {
	return &GlobMuxer{paths: make(map[string]http.Handler)}
}

func (g *GlobMuxer) Handler(req *http.Request) (h http.Handler, pattern string) {
	h, pattern = g.match(req.URL.Path)
	if pattern == "" {
		return http.NotFoundHandler(), pattern
	}
	return
}

func (g *GlobMuxer) Handle(pattern string, h http.Handler) {
	g.paths[pattern] = h
}

func (g *GlobMuxer) match(path string) (h http.Handler, pattern string) {
	longestMatchLen := 0
	for matchPattern, handler := range g.paths {
		if pathMatch(matchPattern, path) {
			matchLen := strings.Count(matchPattern, "/")
			if matchLen > longestMatchLen {
				h = handler
				pattern = matchPattern
			}
		}
	}
	return
}

func pathMatch(pattern, path string) bool {
	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(path, "/")
	for i, patternPart := range patternParts {
		if len(pathParts) > i {
			switch pathParts[i] {
			case "*":
			case patternPart:
			default:
				return false
			}
		}
	}
	return true
}

// example
// g.Handle("/orders/*/products/*")
