package communities

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func debug(r *http.Request, msg string) {
	if log.GetLevel() == log.DebugLevel {
		log.WithFields(log.Fields{
			"url":    r.URL.String(),
			"method": r.Method,
			"id":     r.Header.Get("X-Request-ID"),
		}).Debug(msg)
	}
}
