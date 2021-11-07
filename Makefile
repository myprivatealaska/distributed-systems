GO=$(shell which go)
GO_BUILD=$(GO) build -ldflags="-s -w"


build-server:
