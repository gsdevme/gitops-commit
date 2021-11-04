package slackhttp

import "fmt"

type NamedRepository struct {
	Name       string
	Repository string
	File       string
	Notation   string
	Branch     string
}

type NamedRepositoryRegistry struct {
	r *[]NamedRepository
}

func NewNamedRepositoryRegistry() *NamedRepositoryRegistry {
	return &NamedRepositoryRegistry{
		r: &[]NamedRepository{},
	}
}

func (c *NamedRepositoryRegistry) Add(name string, r string, f string, n string, b string) {
	*c.r = append(*c.r, NamedRepository{
		Name:       name,
		Repository: r,
		File:       f,
		Notation:   n,
		Branch:     b,
	})
}

func (c *NamedRepositoryRegistry) findNamedRepository(n string) (*NamedRepository, error) {
	if c.r == nil {
		return nil, fmt.Errorf("no named repository found for %s", n)
	}

	for _, r := range *c.r {
		if r.Name == n {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no named repository found for %s", n)
}

func (c *NamedRepositoryRegistry) getNamesFlattened() string {
	if c.r == nil {
		return ""
	}

	var names string

	for _, r := range *c.r {
		names += r.Name + ", "
	}

	return names[:len(names)-2]
}
