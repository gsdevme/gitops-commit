package cmd

import (
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
	"os"
)

func NewRootCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "Gitops Commit",
		Short: "A simple tool to mutate a version within a yaml file to trigger gitops operations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	d := os.Getenv("DEBUG")

	if len(d) > 0 {
		log.SetLevel(log.Debug)
	}

	c.AddCommand(newRunCommand())
	c.AddCommand(newServeCommand())

	return &c
}
