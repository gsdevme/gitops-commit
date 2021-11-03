package gitops

import (
	"fmt"
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
	}

	fmt.Println(string(createSimpleYaml("v1.2.99")))

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
