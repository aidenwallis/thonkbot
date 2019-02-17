package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/aidenwallis/thonkbot/mysql"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func (s *Server) logSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeString(w, http.StatusBadRequest, "Please use ?q= parameter!")
		return
	}

	channel := chi.URLParam(r, "channel")

	mins := 5

	minsStr := r.URL.Query().Get("mins")
	if minsStr != "" {
		n, err := strconv.Atoi(minsStr)
		if err == nil && n > 0 {
			mins = n
		}
	}

	quotes, err := mysql.LogSearch(channel, query, mins)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"channel": channel,
			"query":   query,
		}).WithError(err).Error("Failed to get search logs")
		writeString(w, http.StatusInternalServerError, "Failed to get logs from database!")
		return
	}

	usernames := []string{}
	for _, quote := range quotes {
		usernames = append(usernames, quote.Username)
	}

	writeString(w, http.StatusOK, strings.Join(usernames, "\n"))
}
