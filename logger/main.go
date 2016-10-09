// Package logger logs http requests
package logger

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Options ...
type Options struct {
	excludeRoutes []string
	defaultStatus int
}

// WrapWriter ...
type WrapWriter struct {
	http.ResponseWriter
	status int
}

// Write ...
func (ww *WrapWriter) Write(p []byte) (int, error) {
	return ww.ResponseWriter.Write(p)
}

// WriteHeader ...
func (ww *WrapWriter) WriteHeader(status int) {
	ww.status = status
	ww.ResponseWriter.WriteHeader(status)
}

// Wrap ...
func Wrap(h http.Handler, opts *Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if opts == nil {
			opts = &Options{
				excludeRoutes: nil,
				defaultStatus: 200,
			}
		}

		ww := &WrapWriter{ResponseWriter: w, status: opts.defaultStatus}

		start := time.Now()
		h.ServeHTTP(ww, r)
		stop := time.Now()

		for _, i := range opts.excludeRoutes {
			if strings.Contains(r.URL.String(), i) {
				return
			}
		}

		log.Printf("[%d] [%s] %q %v\n", ww.status, r.Method, r.URL.String(), stop.Sub(start))
	}
}
