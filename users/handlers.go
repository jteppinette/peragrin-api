package users

import (
	"net/http"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

func (c *Config) ListHandler(w http.ResponseWriter, r *http.Request) {
	v, err := models.ListUsers(c.Client)
	if err != nil {
		service.Error(w, http.StatusBadRequest, errListUsers)
		return
	}
	json(w, http.StatusBadRequest, v)
}
