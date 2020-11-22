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
	go test . -update -clean
.PHONY: update-golden

lint:
	golangci-lint run --config .golangci.yml
.PHONY: lint

shell:
	cd cmd/dynamosql && \
		AWS_ACCESS_KEY_ID=fake AWS_SECRET_ACCESS_KEY=secret AWS_DEFAULT_REGION=us-west-2 go run *.go --endpoint-url http://localhost:8000
.PHONY: shell