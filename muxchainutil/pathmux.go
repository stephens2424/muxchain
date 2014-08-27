package muxchainutil

import (
	"net/http"
	"strings"
)

// PathMux muxes patterns by globbing over variable components and adding those
// as query parameters to the request for handlers.
type PathMux struct {
	*GlobMux
	patternVariables map[string]map[int]pathVariable // maps pattern component index to variable name
}

// NewPathMuxer initializes a PathMux.
func NewPathMux() *PathMux {
	return &PathMux{NewGlobMux(), make(map[string]map[int]pathVariable)}
}

type pathVariable struct {
	name, def string
}

// Handle registers a handler to a pattern. Patterns may conatain variable components
// specified by a leading colon. For instance, "/order/:id" would map to handlers the
// same as "/order/*" on a GlobMux, however, handlers will see requests as if
// the query were "/order/?id=".
//
// Handlers will also match if a partial query is provided. For instance, /order/x
// will match /order/:id/:name, and the name variable will be empty. Variables are always
// matched from left to right and the handler with the most matches wins (with static strings
// beating variables).
//
// Path components can have default components if they are the last component or are followed
// only by path components with default values. Default components are specified by an extra
// colon, followed by the default value. For example:
//
//    /image/:name/:size:100x100        // valid
//    /image/:name:newimg/:size:100x100 // valid
//    /image/:name/size/:size           // invalid
//    /image/:name:newimg/:size         // invalid
//
// If a passed pattern breaks these rules, this function will panic.
func (p *PathMux) Handle(pattern string, h http.Handler) {
	patternParts := strings.Split(pattern, "/")
	variables := make(map[int]pathVariable)
	foundDefault := false

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			varParts := strings.SplitN(part[1:], ":", 2)
			variable := pathVariable{name: varParts[0]}
			if len(varParts) == 2 {
				foundDefault = true
				variable.def = varParts[1]
			} else if foundDefault {
				panic("pathmux: invalid pattern: cannot have defaulted component followed by non-defaulted component")
			}

			variables[i] = variable
			patternParts[i] = "*"
		} else if foundDefault {
			panic("pathmux: invalid pattern: cannot have defaulted component followed by non-defaulted component")
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
			if variable, ok := variables[i]; ok {
				req.URL.Query().Set(variable.name, part)
				req.Form.Set(variable.name, part)
			}
		}
	}
	h.ServeHTTP(w, req)
}
