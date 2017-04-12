package service

import (
	"net/http"

	"github.com/unrolled/render"
)

var (
	ren = render.New()
)

func Error(w http.ResponseWriter, code int, err error) {
	t := struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		Error  string `json:"error"`
	}{Code: code, Status: http.StatusText(code)}
	if err != nil {
		t.Error = err.Error()
	}
	ren.JSON(w, code, t)
}
