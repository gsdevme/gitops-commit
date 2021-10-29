.PHONY: all
default: build;

fmt:
	go fmt ./...

tests:
	go test ./...

run:
	./dist/gitops-commit_darwin_amd64/gitops-commit

build:
	goreleaser build --snapshot --skip-validate --rm-dist --single-target


all: build run