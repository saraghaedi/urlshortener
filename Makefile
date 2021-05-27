#@IgnoreInspection BashAddShebang

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export APP=urlshortener

export LDFLAGS="-w -s"

all: format lint build

run-server:
	go run -ldflags $(LDFLAGS)  ./cmd/urlshortener server

run-migrate:
	go run -ldflags $(LDFLAGS)  ./cmd/urlshortener migrate

build:
	go build -ldflags $(LDFLAGS)  ./cmd/urlshortener

install:
	go install -ldflags $(LDFLAGS) ./cmd/urlshortener

check-go-bindata:
	which go-bindata || GO111MODULE=off go get -u github.com/jteeuwen/go-bindata/...

bindata: check-go-bindata
	cd internal/app/urlshortener/migrations/postgres && go-bindata -pkg postgres -o ./../bindata/postgres/bindata.go .

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.23.8

lint: check-linter
	golangci-lint -c build/ci/.golangci.yml run $(ROOT)/...

up:
	docker-compose -f ./deployments/docker/urlshortener/docker-compose.yml up -d

down:
	docker-compose -f ./deployments/docker/urlshortener/docker-compose.yml down
