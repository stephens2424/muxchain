package muxchainutil

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stephens2424/muxchain"
)

var DefaultLog *LogHandler

func init() {
	DefaultLog = NewLogHandler(os.Stdout, "", LstdFlags)
}

const (
	Lpath = log.Lshortfile << iota
	Lmethod
	LremoteAddr
	LresponseStatus
	LcontentLength

	LstdFlags = log.Ldate | log.Ltime | Lmethod | Lpath | LresponseStatus | LcontentLength
)

type LogHandler struct {
	flag int
	*log.Logger
}

func NewLogHandler(out io.Writer, prefix string, flag int) *LogHandler {
	return &LogHandler{flag, log.New(out, prefix, flag)}
}

func (l LogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	wr := &responseWriter{ResponseWriter: w}
	l.Println(l.header(wr, req))
}

func (l LogHandler) ServeHTTPChain(w http.ResponseWriter, req *http.Request, h ...http.Handler) {
	wr := &responseWriter{ResponseWriter: w}
	defer func() {
		l.Println(l.header(wr, req))
	}()
	muxchain.HandleChain(wr, req, h...)
}

func (l LogHandler) header(w *responseWriter, req *http.Request) string {
	buf := &bytes.Buffer{}
	if Lmethod&l.flag != 0 {
		fmt.Fprintf(buf, "%s ", req.Method)
	}
	if Lpath&l.flag != 0 {
		fmt.Fprintf(buf, "%s ", req.URL.Path)
	}
	if LresponseStatus&l.flag != 0 {
		fmt.Fprintf(buf, "%d ", w.code)
	}
	if LcontentLength&l.flag != 0 {
		fmt.Fprintf(buf, "%d ", w.written)
	}
	return buf.String()
}

type responseWriter struct {
	written int
	code    int
	http.ResponseWriter
}

func (r *responseWriter) WriteHeader(code int) {
	if r.code == 0 {
		r.code = code
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *responseWriter) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.written += n
	return n, err
}
