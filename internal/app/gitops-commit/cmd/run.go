package cmd

import (
	"fmt"
	"github.com/gsdevme/gitops-commit/internal/pkg/gitops"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func newRunCommand() *cobra.Command {
	c := cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := cmd.Flag("key").Value.String()
			email := cmd.Flag("email").Value.String()
			newVersion := cmd.Flag("version").Value.String()
			notation := cmd.Flag("notation").Value.String()
			repo := strings.TrimRight(cmd.Flag("repo").Value.String(), "/")
			file := strings.TrimLeft(cmd.Flag("file").Value.String(), "/")

			keys, err := gitops.GetPasswordlessKey(key)

			if err != nil {
				return err
			}

			options, c, err := gitops.NewGitOptions(keys)

			if len(email) > 0 {
				options.Email = email
			}

			if err != nil {
				return err
			}

			defer c()

			return gitops.DeployVersionHandler(gitops.DeployVersionCommand{
				GitOptions: *options,
				Repository: repo,
				Notation:   notation,
				File:       file,
				Version:    newVersion,
			})
		},
	}

	c.Flags().String("notation", "", "The yaml path in dot notation i.e. image.tag")
	c.Flags().String("email", "", "The email address of the commit")
	c.Flags().String("version", "", "The semver version you want to deploy i.e. v1.1.2")
	c.Flags().String("key", fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")), "Absolute path to the private key")
	c.Flags().String("repo", "gsdevme/test", "The org/repo path")
	c.Flags().String("file", "/deployments/values.yaml", "The relative path in the repository to the file")

	_ = c.MarkFlagRequired("notation")
	_ = c.MarkFlagRequired("tag")
	_ = c.MarkFlagRequired("file")

	return &c
}
