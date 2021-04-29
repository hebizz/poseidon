GO=go

ver=$(shell cat ./VERSION)
time=$(shell date "+%m/%d/%Y %R %Z")
hash=$(shell git rev-parse --short HEAD)
gover=$(shell go version)

pro=gitlab.jiangxingai.com/poseidon/server

LDFLAG=-ldflags '-X "$(pro)/pkg/config.BuildVersion=$(ver)" \
                 -X "$(pro)/pkg/config.BuildTime=$(time)" \
                 -X "$(pro)/pkg/config.BuildHash=$(hash)" \
                 -X "$(pro)/pkg/config.GoVersion=$(gover)"'

.PHONY: linux-amd64 linux-arm64 windows-amd64 darwin-amd64


linux-amd64: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build  $(LDFLAG) -o bin


linux-arm64: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build  $(LDFLAG) -o bin


windows-amd64: main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build  $(LDFLAG) -o bin

darwin-amd64: main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAG) -o bin



