package slackhttp

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Manifest struct {
	registry *NamedRepositoryRegistry
}

type manifestYaml struct {
	Repositories []struct {
		Name       string `yaml:"name"`
		File       string `yaml:"file"`
		Notation   string `yaml:"notation"`
		Repository string `yaml:"repository"`
		Branch     string `yaml:"branch,omitempty"`
	} `yaml:"repositories"`
}

func (m *Manifest) GetRegistry() *NamedRepositoryRegistry {
	return m.registry
}

func LoadManifest(f string) (*Manifest, error) {
	m := Manifest{
		registry: NewNamedRepositoryRegistry(),
	}

	d, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, fmt.Errorf("cannot read yaml file: %w", err)
	}

	var manifest manifestYaml

	err = yaml.Unmarshal(d, &manifest)

	if err != nil {
		return nil, fmt.Errorf("invalid yaml file: %w", err)
	}

	var branch string

	for _, r := range manifest.Repositories {
		branch = "master"

		if len(r.Branch) > 0 {
			branch = r.Branch
		}

		m.registry.Add(r.Name, r.Repository, r.File, r.Notation, branch)
	}

	return &m, nil
}
