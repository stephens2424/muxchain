package muxchainutil

import (
	"net/http"
	"strings"
)

// PathMux muxes patterns by globbing over variable components and adding those
// as query parameters to the request for handlers.
type PathMux struct {
	*GlobMux
	patternVariables map[string]map[int]string // maps pattern component index to variable name
}

// NewPathMuxer initializes a PathMux.
func NewPathMux() *PathMux {
	return &PathMux{NewGlobMux(), make(map[string]map[int]string)}
}

// Handle registers a handler to a pattern. Patterns may conatain variable components
// specified by a leading colon. For instance, "/order/:id" would map to handlers the
// same as "/order/*" on a GlobMux, however, handlers will see requests as if
// the query were "/order/?id=".
func (p *PathMux) Handle(pattern string, h http.Handler) {
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
	p.GlobMux.Handle(pattern, h)
}

// ServeHTTP handles requests by adding path variables to the request and forwarding
// them to the matched handler. (See Handle)
func (p *PathMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
