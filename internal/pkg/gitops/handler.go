package gitops

import (
	"fmt"
	"io/ioutil"
)

type DeployVersionCommand struct {
	GitOptions GitOptions
	Repository string
	Notation   string
	File       string
	Version    string
}

func DeployVersionHandler(c DeployVersionCommand) error {
	r, err := cloneRepository(&c.GitOptions, c.Repository)

	if err != nil {
		return fmt.Errorf("failed to clone repo %s: %w", c.Repository, err)
	}

	filename := fmt.Sprintf("%s/%s", c.GitOptions.WorkingDirectory, c.File)

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	version, err := ReadCurrentVersion(f, c.Notation)
	if err != nil {
		return fmt.Errorf("cannot read current version deployed: %w", err)
	}

	err = WriteVersion(f, version, c.Version, filename)
	if err != nil {
		return fmt.Errorf("cannot write new version: %w", err)
	}

	return PushVersion(r, &c.GitOptions, c.File, fmt.Sprintf("ci: update tag to %s", c.Version))
}
