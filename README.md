# Gitops Commit

[Under Construction]

A utility designed to mutate a version within a yaml file to trigger
gitops operations i.e. update a helm value

## Usage

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

## TODO

- [ ] Add web service
- [ ] Add webhook to ingress from slack command

## Limitations

- Github only
- YAML only
- Semver tag only