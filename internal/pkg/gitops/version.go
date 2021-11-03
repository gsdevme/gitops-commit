package gitops

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

// ReadCurrentVersion
// Loads the yaml and finds the version at the notation
func ReadCurrentVersion(f []byte, notation string) (string, error) {
	data := make(map[string]interface{})
	err := yaml.Unmarshal(f, &data)

	if err != nil || len(data) <= 0 {
		return "", fmt.Errorf("unvalid valid: %w", err)
	}

	v, err := unwrapYaml(data, notation)

	if err != nil {
		return "", fmt.Errorf("cannot find current version within yaml at %s: %w", notation, err)
	}

	return v, nil
}

func WriteVersion(f []byte, version string, newVersion string, filename string) error {
	output := bytes.Replace(f, []byte(version), []byte(newVersion), -1)
	err := ioutil.WriteFile(filename, output, 0666)

	if err != nil {
		return fmt.Errorf("cannot replace version: %w", err)
	}

	return nil
}

func unwrapYaml(yaml map[string]interface{}, notion string) (string, error) {
	d := yaml

	for _, k := range strings.Split(notion, ".") {
		if _, ok := d[k]; !ok {
			return "", fmt.Errorf("unable to find %s in yaml", notion)
		}

		if v, ok := d[k].(string); ok {
			return v, nil
		}

		if m, ok := d[k].(map[string]interface{}); ok {
			for key, value := range m {
				d[key] = value
			}
		}
	}

	return "", errors.New("unable to find in yaml")
}
