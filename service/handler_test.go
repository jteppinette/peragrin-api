package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		handler Handler
		code    int
		bytes   []byte
	}{
		{
			func(r *http.Request) *Response {
				data := struct {
					Detail string `json:"detail"`
				}{"hey"}
				return NewResponse(nil, http.StatusOK, data)
			},
			http.StatusOK,
			[]byte(`{"detail":"hey"}`),
		},
		{
			func(r *http.Request) *Response {
				return NewResponse(nil, http.StatusOK, nil)
			},
			http.StatusOK,
			nil,
		},
		{
			func(r *http.Request) *Response {
				return nil
			},
			http.StatusOK,
			nil,
		},
	}
	for _, test := range tests {
		w := httptest.NewRecorder()
		test.handler.ServeHTTP(w, &http.Request{})
		if w.Code != test.code {
			t.Errorf("expected code to be %d, got %d", test.code, w.Code)
		}
		b := w.Body.Bytes()
		if !bytes.Equal(b, test.bytes) {
			t.Errorf("expected response to be %s, got %s", test.bytes, b)
		}
	}
}
