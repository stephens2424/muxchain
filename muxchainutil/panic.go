package muxchainutil

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/stephens2424/muxchain"
)

// DefaultPanicRecovery is a handler that enables basic panic recovery
// for all handlers chained after it.
var DefaultPanicRecovery = PanicRecovery{http.HandlerFunc(DefaultRecoverFunc)}

func DefaultRecoverFunc(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

type PanicRecovery struct {
	Recovery http.Handler
}

func (p PanicRecovery) ServeHTTPChain(w http.ResponseWriter, req *http.Request, h ...http.Handler) {
	recovery := p.Recovery
	if recovery == nil {
		recovery = http.HandlerFunc(DefaultRecoverFunc)
	}
	defer func() {
		if r := recover(); r != nil {
			recovery.ServeHTTP(w, req)
			fmt.Fprintln(os.Stderr, r)
			debug.PrintStack()
		}
	}()
	muxchain.HandleChain(w, req, h...)
}

func (p PanicRecovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.ServeHTTPChain(w, req, nil)
}
