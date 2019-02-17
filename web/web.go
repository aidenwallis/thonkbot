package web

import (
	"net/http"
	"time"

	"github.com/aidenwallis/thonkbot/botmanager"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	webhost string
	manager *botmanager.BotManager
	log     logrus.FieldLogger
	Router  *chi.Mux
}

func New(webhost string, manager *botmanager.BotManager, logger logrus.FieldLogger) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	s := &Server{
		webhost: webhost,
		manager: manager,
		Router:  r,
		log:     logger.WithField("package", "web"),
	}

	s.mountRoutes(r)
	return s
}

func (s *Server) Start() {
	s.log.Info("Listening on ", s.webhost)
	http.ListenAndServe(s.webhost, s.Router)
}

func writeString(w http.ResponseWriter, code int, s string) {
	write(w, code, []byte(s))
}

func write(w http.ResponseWriter, code int, bs []byte) {
	w.WriteHeader(code)
	_, _ = w.Write(bs)
}
