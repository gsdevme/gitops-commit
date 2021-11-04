package slackhttp

import "testing"

func TestNamedRepositoryRegistry_getNamesFlattened(t *testing.T) {
	type fields struct {
		r *[]NamedRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "gg",
			fields: fields{
				r: &[]NamedRepository{
					{
						Name:       "repo1",
						Repository: "test",
						File:       "test",
						Notation:   "test",
						Branch:     "master",
					},
					{
						Name:       "repo2",
						Repository: "test",
						File:       "test",
						Notation:   "test",
						Branch:     "master",
					},
				},
			},
			want: "repo1, repo2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &NamedRepositoryRegistry{
				r: tt.fields.r,
			}
			if got := c.getNamesFlattened(); got != tt.want {
				t.Errorf("getNamesFlattened() = %v, want %v", got, tt.want)
			}
		})
	}
}
