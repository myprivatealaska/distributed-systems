.DEFAULT_GOAL := build

GO=$(shell which go)
GO_BUILD=$(GO) build -ldflags="-s -w"
GO_STATIC_BUILD=GOARCH=amd64 GOOS=linux CGO_ENABLED=0 $(GO) build -installsuffix 'static'

build: modules
	$(GO_BUILD) -o bin/server .

modules:
	$(GO) mod download
