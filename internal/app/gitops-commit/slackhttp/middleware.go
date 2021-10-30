package slackhttp

import (
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"net/http"
)

func (s *server) SlackCommandMiddleware(next func(w http.ResponseWriter, s slack.SlashCommand)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, s.slackSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		next(w, s)
	}
}
