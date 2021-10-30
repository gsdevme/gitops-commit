package slackhttp

import (
	"encoding/json"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type server struct {
	router      *mux.Router
	slackSecret string
	keys        *ssh.PublicKeys
	registry    *NamedRepositoryRegistry
}

func NewSlackCommandServer(r NamedRepositoryRegistry, keys *ssh.PublicKeys) *server {
	s := &server{
		router:      mux.NewRouter(),
		keys:        keys,
		slackSecret: os.Getenv("SLACK_SIGNING_SECRET"),
		registry:    &r,
	}

	s.routes()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/ping" {
			body, err := json.Marshal(struct {
				Pong bool `json:"pong"`
			}{
				Pong: true,
			})

			if err != nil {
				return
			}

			s.respond(w, r, body, http.StatusOK)

			return
		}

		body, err := json.Marshal(struct {
			Ok bool `json:"ok"`
		}{
			Ok: true,
		})

		if err != nil {
			return
		}

		s.respond(w, r, body, http.StatusOK)
	}
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("server", "gitops-commit")
	w.Header().Set("content-type", "text/json")
	w.WriteHeader(status)

	if response, ok := data.([]byte); ok {
		_, err := w.Write(response)
		if err != nil {
			return
		}

		return
	}

	if response, ok := data.(string); ok {
		_, err := w.Write([]byte(response))
		if err != nil {
			return
		}

		return

	}
}

func (s *server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
}
