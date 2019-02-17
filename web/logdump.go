package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aidenwallis/thonkbot/mysql"
	"github.com/go-chi/chi"
)

func (s *Server) getLogdump(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")

	quotes, err := mysql.GetAllLogs(channel)
	if err != nil {
		s.log.WithField("channel", channel).WithError(err).Error("Error while fetching logdump")
		writeString(w, http.StatusInternalServerError, "Database error while fetching logdump.")
		return
	}

	quoteString := []string{}
	for _, quote := range quotes {
		formattedDate := quote.CreatedAt.UTC().Format(time.ANSIC)
		quoteString = append(quoteString, fmt.Sprintf("[%s] %s: %s", formattedDate, quote.Username, quote.Message))
	}

	writeString(w, http.StatusOK, strings.Join(quoteString, "\n"))
}
