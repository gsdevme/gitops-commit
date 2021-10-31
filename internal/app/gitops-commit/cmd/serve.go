package cmd

import (
	"fmt"
	"github.com/google/martian/log"
	"github.com/gsdevme/gitops-commit/internal/app/gitops-commit/slackhttp"
	"github.com/gsdevme/gitops-commit/internal/pkg/gitops"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func newServeCommand() *cobra.Command {
	c := cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			port := cmd.Flag("port").Value.String()
			key := cmd.Flag("key").Value.String()
			file := cmd.Flag("manifest").Value.String()

			keys, err := gitops.GetPasswordlessKey(key)

			if err != nil {
				return err
			}

			manifest, err := slackhttp.LoadManifest(file)
			if err != nil {
				return err
			}

			r := manifest.GetRegistry()

			s := slackhttp.NewSlackCommandServer(*r, keys)
			server := &http.Server{
				Addr:         fmt.Sprintf(":%s", port),
				Handler:      s,
				ReadTimeout:  3 * time.Second,
				WriteTimeout: 5 * time.Second,
				ConnState: func(conn net.Conn, state http.ConnState) {
					log.Debugf("%s - %s\n", conn.RemoteAddr(), state.String())
				},
			}

			fmt.Fprintf(os.Stdout, "Listening on http://127.0.0.1:%s", port)

			defer server.Close()

			return server.ListenAndServe()
		},
	}

	p := 8080

	if len(os.Getenv("PORT")) > 0 {
		envPort, err := strconv.Atoi(os.Getenv("PORT"))

		if err != nil {
			p = envPort
		}
	}

	c.Flags().Int("port", p, "The web server port to listen on")
	c.Flags().String("key", fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")), "Absolute path to the private key")
	c.Flags().String("manifest", fmt.Sprintf("%s/.gitops-commit/manifest.yaml", os.Getenv("HOME")), "Absolute path to the manifest")

	return &c
}
