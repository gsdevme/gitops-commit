package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "Gitops Commit",
		Short: "A simple tool to mutate a version within a yaml file to trigger gitops operations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	c.AddCommand(newRunCommand())

	return &c
}
