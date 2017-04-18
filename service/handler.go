package service

import (
	"net/http"

	"github.com/unrolled/render"
)

var (
	rend = render.New().JSON
)

// Response represents the return value to the Handler type.
// This struct encapsulates the information need to log and write
// at the end of a request/response cycle.
type Response struct {
	Error error
	Code  int
	Data  interface{}
}

// NewResponse returns an initialized response pointer.
func NewResponse(err error, code int, data interface{}) *Response {
	return &Response{err, code, data}
}

// Handler overrides the typical http.Handler interface with the Response
// return value. This allows the types ServeHTTP function to handle the Response
// in a single place providing standardized logging and response writing.
type Handler func(r *http.Request) *Response

// ServeHTTP calls the handler's underlying function, and it will properly
// handle the returned response.
// If the response or response data is nil, an empty text/html response will
// be returned. Otherwise, the response data will be written as encoded JSON.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h(r)

	// If the response is nil, then use http.StatusOK as the default HTTP code.
	var code int
	if response == nil {
		code = http.StatusOK
	} else {
		code = response.Code
	}

	// If the response or response data is nil, then write the calculated code
	// as a default empty text/html response.
	if response == nil || response.Data == nil {
		w.WriteHeader(code)
		return
	}

	rend(w, code, response.Data)
}
