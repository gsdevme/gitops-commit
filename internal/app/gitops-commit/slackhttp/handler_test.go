package slackhttp

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleSlackCommand(t *testing.T) {
	t.Run("invalid command", func(t *testing.T) {
		s := &server{}
		r := &NamedRepositoryRegistry{}

		tests := []struct{
			name string
			command string
			expect string
		}{
			{
				name: "when using missing bits",
				command: "totally wrong",
				expect: "Incorrect usage, expected /gitops-commit [command] [name] [tag]",
			},
			{
				name: "when using an invalid command",
				command: "wrong thing v1.2.3",
				expect: "Unknown command 'wrong', expected /gitops-commit [command] [name] [tag]",
			},
			{
				name: "when using a valid command ",
				command: "deploy thing v1.2.3",
				expect: "Unknown named repository, cannot handle \"thing\", availabe options ()",
			},
		}

		for _, tt := range tests {
			slackCommand := slack.SlashCommand{
				Text: tt.command,
			}

			t.Run(tt.name, func(t *testing.T) {
				rr := httptest.NewRecorder()

				s.handleSlackCommand(r)(rr, slackCommand)

				if c := rr.Code; c != http.StatusOK {
					t.Errorf("handleSlackCommand() = %v, want %v", c, http.StatusOK)
				}

				response := make(map[string]interface{})
				err := json.Unmarshal(rr.Body.Bytes(), &response)

				if err != nil {
					t.Errorf("handleSlackCommand() invalid json: %s", err)
				}

				if got, ok := response["text"].(string); ok {
					if got != tt.expect {
						t.Errorf("expected '%s', got '%s'", tt.expect, got)
					}
				}
			})
		}
	})
}