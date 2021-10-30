package slackhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gsdevme/gitops-commit/internal/pkg/gitops"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

func (s *server) handleSlackCommand(registry *NamedRepositoryRegistry) func(http.ResponseWriter, slack.SlashCommand) {
	return func(w http.ResponseWriter, sl slack.SlashCommand) {
		text := strings.Split(sl.Text, " ")

		if len(text) < 3 || len(text) > 3 {
			sendEphemeralMsg("Incorrect usage, expected /gitops-commit [command] [name] [tag]", w)

			return
		}

		switch text[0] {
		case "deploy":
			deploy(s, w, registry, text[1], text[2])

			return
		}
	}
}

func deploy(s *server, w http.ResponseWriter, registry *NamedRepositoryRegistry, name string, version string) {

	r, err := registry.findNamedRepository(name)

	if err != nil {
		sendEphemeralMsg(fmt.Sprintf("Unknown named repository, cannot handle \"%s\", availabe options (%s)", name, registry.getNamesFlattened()), w)

		return
	}

	if len(version) != 7 {
		sendEphemeralMsg(fmt.Sprintf("version does not look semver? %s", version), w)

		return
	}

	options, f, err := gitops.NewGitOptions(s.keys)
	if err != nil {
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
			sendEphemeralMsg(fmt.Sprintf("failed to deploy: %s", err), w)

			return
		}
	}()

	params := &slack.Msg{
		Text:         fmt.Sprintf(":alert: Deploying tag `%s` to `%s`:`%s`", version, r.Repository, r.File),
		ResponseType: slack.ResponseTypeInChannel,
	}
	b, err := json.Marshal(params)
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

func sendEphemeralMsg(m string, w http.ResponseWriter) {
	b, err := json.Marshal(slack.Msg{
		Text:         m,
		ResponseType: slack.ResponseTypeEphemeral,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		return
	}
}
