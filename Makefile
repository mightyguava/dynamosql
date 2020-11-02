SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

dynamodb:
	docker run --name dynamodb -d -p 8000:8000 amazon/dynamodb-local
.PHONY: dynamodb

test:
	go test -p 1 ./...

update-golden:
	go test ./parser -update -clean
	go test ./querybuilder -update -clean
	go test ./driver -update -clean

lint:
	golangci-lint run --config .golangci.yml