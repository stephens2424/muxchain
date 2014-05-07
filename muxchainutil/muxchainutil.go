package muxchainutil

import (
	"net/http"

	"stephensearles.com/muxchain"
)

// Default is a handler that enables panic recovery, logging to standard out, and gzip for all
// request paths chained after it.
var Default http.Handler

func init() {
	m := &muxchain.MuxChain{}
	m.Chain("/", DefaultPanicRecovery, DefaultLog, Gzip)
	Default = m
}
