# Gitops Commit

[Under Construction]

A utility designed to mutate a version within a yaml file to trigger
gitops operations i.e. update a helm value

## Usage

Within the `~/.gitops-commit/manifest.yaml`

```yaml
repositories:
  # Example of what the fields map to
  # https://github.com/gsdevme/test/blob/master/deployments/foo/values.yaml
  # https://github.com/gsdevme/test
  - name: my-app-test
    repository: gsdevme/test
    file: deployments/test/values.yaml
    notation: image.tag
  - name: my-app-test
    repository: gsdevme/test
    file: deployments/prod/values.yaml
    notation: image.tag
    branch: master
```

---

```bash
$: gitops-commit  -h                                                                                    
Usage:
  gitops-commit [flags]

Flags:
      --email string     The email address of the commit
      --file string      The relative path in the repository to the file (default "/deployments/values.yaml")
  -h, --help             help for gitops-commit
      --key string       Absolute path to the private key (default "/Users/gavin/.ssh/id_rsa")
      --repo string      The org/repo path (default "gsdevme/test")
      --version string   The semver version you want to deploy i.e. v1.1.2

```

## Limitations

- Github only
- YAML only
- Passwordless keys only
