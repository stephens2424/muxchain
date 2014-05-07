package muxchainutil

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var DefaultLog *LogHandler

func init() {
	DefaultLog = NewLogHandler(os.Stdout, "", LstdFlags)
}

const (
	Lpath = log.Lshortfile << iota
	Lmethod
	LremoteAddr

	LstdFlags = log.Ldate | log.Ltime | Lmethod | Lpath
)

type LogHandler struct {
	flag int
	*log.Logger
}

func NewLogHandler(out io.Writer, prefix string, flag int) *LogHandler {
	return &LogHandler{flag, log.New(out, prefix, flag)}
}

func (l LogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	l.Println(l.header(req))
}

func (l LogHandler) header(req *http.Request) string {
	buf := &bytes.Buffer{}
	if Lmethod&l.flag != 0 {
		fmt.Fprintf(buf, "%s ", req.Method)
	}
	if Lpath&l.flag != 0 {
		fmt.Fprintf(buf, "%s ", req.URL.Path)
	}
	return buf.String()
}
