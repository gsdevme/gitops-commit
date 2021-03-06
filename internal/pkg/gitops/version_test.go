package gitops

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestReadCurrentVersion(t *testing.T) {
	type args struct {
		f        []byte
		notation string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty file",
			args: args{
				f:        []byte(""),
				notation: "test.test",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid yaml",
			args: args{
				f:        []byte("wibble&&3..\nfff"),
				notation: "image.tag",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "simple yaml",
			args: args{
				f:        createSimpleYaml("v1.2.99"),
				notation: "image.tag",
			},
			want:    "v1.2.99",
			wantErr: false,
		},
		{
			name: "invalid notation to yaml",
			args: args{
				f:        createSimpleYaml("v1.2.99"),
				notation: "image.nope",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "array notation",
			args: args{
				f:        createSimpleArrayYaml("v1.9"),
				notation: "images[3].tag",
			},
			want:    "v1.9",
			wantErr: false,
		},
		{
			name: "helm example",
			args: args{
				f:        createHelmYaml("v1.9"),
				notation: "dependencies[0].version",
			},
			want:    "v1.9",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCurrentVersion(tt.args.f, tt.args.notation)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCurrentVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadCurrentVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createSimpleYaml(t string) []byte {
	type tag struct {
		Tag string `yaml:"tag"`
	}

	marshal, err := yaml.Marshal(&struct {
		Image tag `yaml:"image"`
	}{
		Image: tag{Tag: t},
	})
	if err != nil {
		return nil
	}

	return marshal
}

func createSimpleArrayYaml(t string) []byte {
	type tag struct {
		Tag string `yaml:"tag"`
	}

	marshal, err := yaml.Marshal(&struct {
		Images []tag `yaml:"images"`
	}{
		Images: []tag{{Tag: "v1"}, {Tag: "v2"}, {Tag: "v3"}, {Tag: t}},
	})
	if err != nil {
		return nil
	}

	return marshal
}

func createHelmYaml(t string) []byte {
	type dependencies struct {
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Repository string `yaml:"repository"`
	}

	marshal, err := yaml.Marshal(&struct {
		Dependencies []dependencies `yaml:"dependencies"`
	}{
		Dependencies: []dependencies{{
			Name:       "example-test",
			Version:    t,
			Repository: "gsdevme/test",
		}},
	})
	if err != nil {
		return nil
	}

	return marshal
}
