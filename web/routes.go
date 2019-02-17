package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) mountRoutes(r chi.Router) {
	r.Get("/", healthStatus)

	r.Get("/{channel}/logs/{username}", s.getChannelLogs)
	r.Get("/{channel}/logdump", s.getLogdump)
	r.Get("/{channel}/search", s.logSearch)
}

func healthStatus(w http.ResponseWriter, r *http.Request) {
	writeString(w, http.StatusOK, "Welcome to Thonkbot!")
}
