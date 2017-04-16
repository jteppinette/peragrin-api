package users

import (
	"net/http"

	"gitlab.com/peragrin/api/models"
)

// ListHandler writes all users in the database to the provided response writer.
func (c *Config) ListHandler(w http.ResponseWriter, r *http.Request) {
	v, err := models.ListUsers(c.Client)
	if err != nil {
		// errListUsers
		rend(w, http.StatusBadRequest, nil)
		return
	}
	rend(w, http.StatusBadRequest, v)
}
