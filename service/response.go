package service

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
