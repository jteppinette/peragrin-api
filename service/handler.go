package service

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

// Handler overrides the typical http.Handler interface with the Response
// return value. This allows the types ServeHTTP function to handle the Response
// in a single place providing standardized logging and response writing.
type Handler func(r *http.Request) *Response

// ServeHTTP calls the handler's underlying function, and it will properly
// handle the returned response.
// If the response or response data is nil, an empty text/html response will
// be returned. Otherwise, the response data will be written as encoded JSON.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"stack": strings.Replace(strings.Replace(string(debug.Stack()), "\n", " || ", -1), "\t", "", -1),
				"error": err,
			}).Error("panic")
		}
	}()

	start := time.Now()

	if log.GetLevel() == log.DebugLevel {
		b, _ := httputil.DumpRequest(r, true)
		log.WithFields(log.Fields{
			"request": strings.Replace(string(b), "\r\n", " ", -1),
		}).Debug("request dump")
	}

	response := h(r)

	// If the response is nil, then use http.StatusOK as the default HTTP code.
	var code int
	if response == nil {
		code = http.StatusOK
	} else {
		code = response.Code
	}

	// Log the the outgoing request/response.
	fields := log.Fields{
		"code":     code,
		"method":   r.Method,
		"id":       r.Header.Get("X-Request-ID"),
		"delta-ns": time.Now().Sub(start).Nanoseconds(),
	}
	if r.URL != nil {
		fields["url"] = r.URL.String()
	}
	if log.GetLevel() == log.DebugLevel {
		fields["response-data"] = fmt.Sprintf("%+v", response.Data)
	}
	if ip := getIPAddress(r); ip != "" {
		fields["ip"] = ip
	}
	if response != nil && response.Error != nil {
		fields["err"] = response.Error
		log.WithFields(fields).Error("access log")
	} else {
		log.WithFields(fields).Info("access log")
	}

	// If the response or response data is nil, then write the calculated code
	// as a default empty text/html response.
	if response == nil || response.Data == nil {
		w.WriteHeader(code)
		return
	}

	rend(w, code, response.Data)
}
