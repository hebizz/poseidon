GO=go

ver=$(shell cat ./VERSION)
time=$(shell date "+%m/%d/%Y %R %Z")
hash=$(shell git rev-parse --short HEAD)
gover=$(shell go version)

pro=gitlab.jiangxingai.com/poseidon/client

LDFLAG=-ldflags '-X "$(pro)/pkg/config.BuildVersion=$(ver)" \
                 -X "$(pro)/pkg/config.BuildTime=$(time)" \
                 -X "$(pro)/pkg/config.BuildHash=$(hash)" \
                 -X "$(pro)/pkg/config.GoVersion=$(gover)" \
                 -X "$(pro)/pkg/config.GoBuildType=amd" '

LDFLAG_ARM=-ldflags '-X "$(pro)/pkg/config.BuildVersion=$(ver)" \
			  -X "$(pro)/pkg/config.BuildTime=$(time)" \
			  -X "$(pro)/pkg/config.BuildHash=$(hash)" \
			  -X "$(pro)/pkg/config.GoVersion=$(gover)" \
			  -X "$(pro)/pkg/config.GoBuildType=arm" '

.PHONY: linux-amd64 linux-arm64 linux-armV7 windows-amd64 darwin-amd64 macos-amd64


linux-amd64: main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build  $(LDFLAG) -o webConfig-amd

#todo:指定编译器在交叉编译cgo时仍然会报错
linux-arm64: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC=arm-linux-gnueabi-gcc $(GO) build  $(LDFLAG_ARM) -o webConfig-arm

linux-armV7: main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabi-gcc  $(GO) build  $(LDFLAG_ARM) -o webConfig-armV7

windows-amd64: main.go
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GO) build  $(LDFLAG) -o bin
	# brew install mingw-w64
	# CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -v

macos-amd64: main.go
	#brew install FiloSottile/musl-cross/musl-cross
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build $(LDFLAG) -a -v -o webConfig-amd

macos-arm64: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64  $(GO) build  $(LDFLAG_ARM) -o webConfig-arm

