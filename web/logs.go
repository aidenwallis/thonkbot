package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aidenwallis/thonkbot/mysql"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func (s *Server) getChannelDump(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) getChannelLogs(w http.ResponseWriter, r *http.Request) {
	limit := 100

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err == nil && n > 0 {
			limit = n
		}
	}

	channel := chi.URLParam(r, "channel")
	username := chi.URLParam(r, "username")

	quotes, err := mysql.GetUserLogs(channel, username, limit)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"channel":  channel,
			"username": username,
		}).WithError(err).Error("Failed to get logs")
		writeString(w, http.StatusInternalServerError, "Failed to get logs from database!")
		return
	}

	quoteString := []string{}
	for _, quote := range quotes {
		formattedDate := quote.CreatedAt.UTC().Format(time.ANSIC)
		quoteString = append(quoteString, fmt.Sprintf("[%s] %s: %s", formattedDate, quote.Username, quote.Message))
	}

	writeString(w, http.StatusOK, strings.Join(quoteString, "\n"))
}
