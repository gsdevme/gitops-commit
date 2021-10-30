package gitops

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io/ioutil"
	"os"
	"time"
)

type GitOptions struct {
	WorkingDirectory string
	Keys             ssh.PublicKeys
	Branch           string
	Email            string
	Name             string
}

func NewGitOptions(key string) (*GitOptions, func(), error) {
	dir, err := ioutil.TempDir("/tmp", "prefix")
	if err != nil {
		return nil, nil, err
	}

	keys, err := getPasswordlessKey(key)

	if err != nil {
		return nil, nil, err
	}

	return &GitOptions{
			WorkingDirectory: dir,
			Keys:             *keys,
			Branch:           "master",
			Name:             "gitops-commit-bot",
			Email:            "gitops-commit@example.com",
		}, func() {
			err := os.RemoveAll(dir)
			if err != nil {
				return
			}
		}, nil
}

func PushVersion(r *git.Repository, options *GitOptions, file string, message string) {
	tree, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	_, err = tree.Add(file)

	if err != nil {
		panic(err)
	}

	commit, err := tree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  options.Name,
			Email: options.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		panic(err)
	}

	_, err = r.CommitObject(commit)

	if err != nil {
		panic(err)
	}

	err = r.Push(&git.PushOptions{
		Auth: &options.Keys,
	})

	if err != nil {
		panic(err)
	}
}

func getPasswordlessKey(key string) (*ssh.PublicKeys, error) {
	publicKeys, err := ssh.NewPublicKeysFromFile("git", key, "")
	if err != nil {
		return nil, fmt.Errorf("private/public key invalid: %w", err)
	}

	return publicKeys, nil
}

func cloneRepository(o *GitOptions, r string) (*git.Repository, error) {

	return git.PlainClone(o.WorkingDirectory, false, &git.CloneOptions{
		Auth:         &o.Keys,
		URL:          fmt.Sprintf("git@github.com:%s.git", r),
		SingleBranch: true,
		Depth:        1,
	})
}
