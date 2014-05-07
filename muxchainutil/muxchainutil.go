package muxchainutil

import (
	"net/http"

	"stephensearles.com/muxchain"
)

var Default http.Handler

func init() {
	m := &muxchain.MuxChain{}
	m.Chain("/", DefaultPanicRecovery, DefaultLog, Gzip)
	Default = m
}
