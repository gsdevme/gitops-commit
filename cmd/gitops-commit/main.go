package main

import (
	"fmt"
	"github.com/gsdevme/gitops-commit/internal/app/gitops-commit/cmd"
	"os"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
