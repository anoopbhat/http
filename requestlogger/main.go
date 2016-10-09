// Package requestlogger ...
package requestlogger

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Options ...
type Options struct {
	excludeRoutes []string
}

// LogRecord ...
type LogRecord struct {
	http.ResponseWriter
	status int
}

// Write ...
func (r *LogRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

// WriteHeader ...
func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// WrapHandler ...
func WrapHandler(h http.Handler, opts *Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := &LogRecord{ResponseWriter: w, status: 200}

		start := time.Now()
		h.ServeHTTP(lw, r)
		stop := time.Now()

		if opts == nil {
			opts = &Options{excludeRoutes: nil}
		}

		for _, i := range opts.excludeRoutes {
			if strings.Contains(r.URL.String(), i) {
				return
			}
		}

		log.Printf("[%d] [%s] %q %v\n", lw.status, r.Method, r.URL.String(), stop.Sub(start))
	}
}
