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
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			key := cmd.Flag("key").Value.String()
			repo := strings.TrimRight(cmd.Flag("repo").Value.String(), "/")
			file := strings.TrimLeft(cmd.Flag("file").Value.String(), "/")

			options, c, err := gitops.NewGitOptions(key)

			if err != nil {
				panic(err.Error()) // todo
			}

			defer c()

			r, err := gitops.CloneRepository(options, repo)

			if err != nil {
				panic(err.Error()) // todo
			}

			filename := fmt.Sprintf("%s/%s", options.WorkingDirectory, file)

			f, err := ioutil.ReadFile(filename)

			if err != nil {
				panic(fmt.Errorf("cannot read file: %w", err))
			}

			version, err := gitops.ReadCurrentVersion(f, "foo.woo.wibble.tag")
			if err != nil {
				panic(err) // todo
			}

			err = gitops.WriteVersion(f, version, "v2.3.1", filename)
			if err != nil {
				panic(err) // todo
			}

			gitops.PushVersion(r, options, file)
		},
	}

	c.Flags().String("key", fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")), "Absolute path to the private key")
	c.Flags().String("repo", "gsdevme/test", "The org/repo path")
	c.Flags().String("file", "/deployments/values.yaml", "The relative path in the repository to the file")

	_ = c.MarkFlagRequired("file")

	return &c
}

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
