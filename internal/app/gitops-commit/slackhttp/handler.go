package slackhttp

import (
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"github.com/gsdevme/gitops-commit/internal/pkg/gitops"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

const (
	SlackUnknownResponse        = "Incorrect usage, expected /gitops-commit [command] [name] [tag]"
	SlackUnknownCommandResponse = "Unknown command '%s', expected /gitops-commit [command] [name] [tag]"
)

func (s *server) handleSlackCommand(registry *NamedRepositoryRegistry) func(http.ResponseWriter, slack.SlashCommand) {
	return func(w http.ResponseWriter, sl slack.SlashCommand) {
		text := strings.Split(sl.Text, " ")

		if len(text) < 3 || len(text) > 3 {
			respondSlack(SlackUnknownResponse, slack.ResponseTypeEphemeral, w)

			return
		}

		command := text[0]

		switch command {
		case "deploy":
			deploy(s, w, registry, text[1], text[2])

			return
		default:
			respondSlack(
				fmt.Sprintf(SlackUnknownCommandResponse, command),
				slack.ResponseTypeEphemeral,
				w,
			)
		}
	}
}

func deploy(s *server, w http.ResponseWriter, registry *NamedRepositoryRegistry, name string, version string) {
	r, err := registry.findNamedRepository(name)

	if err != nil {
		respondSlack(
			fmt.Sprintf("Unknown named repository, cannot handle \"%s\", availabe options (%s)", name, registry.getNamesFlattened()),
			slack.ResponseTypeEphemeral,
			w,
		)

		return
	}

	if len(version) != 7 {
		respondSlack(fmt.Sprintf("version does not look semver? %s", version), slack.ResponseTypeEphemeral, w)

		return
	}

	options, f, err := gitops.NewGitOptions(s.keys)
	if err != nil {
		log.Errorf("failed to setup environment: %w", err)

		return
	}

	defer f()

	command := gitops.DeployVersionCommand{
		GitOptions: *options,
		Repository: r.Repository,
		Notation:   r.Notation,
		File:       r.File,
		Version:    version,
	}

	go func() {
		err = gitops.DeployVersionHandler(command)
		if err != nil {
			log.Errorf("failed to deploy %w", err)

			return
		}
	}()

	respondSlack(
		fmt.Sprintf(":alert: Deploying tag `%s` to `%s`:`%s`", version, r.Repository, r.File),
		slack.ResponseTypeInChannel,
		w,
	)
}

func respondSlack(m string, t string, w http.ResponseWriter) {
	b, err := json.Marshal(slack.Msg{
		Text:         m,
		ResponseType: t,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		return
	}
}
