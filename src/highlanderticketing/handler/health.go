package handler

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Infof("healthy")
	io.WriteString(w, `{"alive": true}`)
}
