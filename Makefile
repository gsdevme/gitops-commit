.PHONY: all
default: build;

fmt:
	go fmt ./...

lint:
	 golangci-lint run

tests:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

run:
	./dist/gitops-commit_darwin_amd64/gitops-commit

build:
	goreleaser build --snapshot --skip-validate --rm-dist --single-target


all: build run