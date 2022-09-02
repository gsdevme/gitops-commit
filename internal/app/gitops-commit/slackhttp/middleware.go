package slackhttp

import (
	"github.com/google/martian/log"
	"github.com/slack-go/slack"
	"io"
	"net/http"
)

func (s *server) SlackCommandMiddleware(next func(w http.ResponseWriter, s slack.SlashCommand)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, s.slackSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			log.Debugf("unauthorized request")

			return
		}

		next(w, s)
	}
}
