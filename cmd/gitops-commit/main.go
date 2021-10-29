package main

import (
	"fmt"
	"github.com/gsdevme/gitops-commit/internal/pkg/gitops"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func NewRootCommand() *cobra.Command {
	c := cobra.Command{
		Use: "gitops-commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := cmd.Flag("key").Value.String()
			email := cmd.Flag("email").Value.String()
			newVersion := cmd.Flag("version").Value.String()
			notation := cmd.Flag("notation").Value.String()
			repo := strings.TrimRight(cmd.Flag("repo").Value.String(), "/")
			file := strings.TrimLeft(cmd.Flag("file").Value.String(), "/")

			options, c, err := gitops.NewGitOptions(key)

			if len(email) > 0 {
				options.Email = email
			}

			if err != nil {
				return err
			}

			defer c()

			r, err := gitops.CloneRepository(options, repo)

			if err != nil {
				return fmt.Errorf("failed to clone repo %s: %w", repo, err)
			}

			filename := fmt.Sprintf("%s/%s", options.WorkingDirectory, file)

			f, err := ioutil.ReadFile(filename)

			if err != nil {
				return fmt.Errorf("cannot read file: %w", err)
			}

			version, err := gitops.ReadCurrentVersion(f, notation)
			if err != nil {
				return fmt.Errorf("cannot read current version deployed: %w", err)
			}

			err = gitops.WriteVersion(f, version, newVersion, filename)
			if err != nil {
				return fmt.Errorf("cannot write new version: %w", err)
			}

			gitops.PushVersion(r, options, file, fmt.Sprintf("ci: update tag to %s", newVersion))

			return nil
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

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
