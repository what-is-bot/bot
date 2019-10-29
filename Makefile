APP_NAME = what-is/bot
VERSION = $(shell git rev-parse --short HEAD)
TEMP_OS_VERSION = $(shell uname)
OS_VERSION = $(shell echo $(TEMP_OS_VERSION) | tr '[:upper:]' '[:lower:]')
USR_ID = $(shell id -u $(whoami))
GRP_ID = $(shell id -g $(whoami))

DOCKER_CMD=docker run \
		--rm -it \
		-v "$(shell pwd)":"/usr/src/$(APP_NAME)" \
		-w "/usr/src/$(APP_NAME)" \
		-u root

DOCKER_IMG=$(DOCKER_CMD) golang:1.13.3

FIX_PERM=chown -R $(USR_ID):$(GRP_ID) ./out/

GO_GET_CMD = go get -v ./...
GO_TEST_CMD = go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
GO_BUILD_WIN_CMD = CGO_ENABLED=0 GOOS=windows go build -v -a --ldflags="-s" -o ./out/bot.exe ./cmd/bot
GO_BUILD_CMD = CGO_ENABLED=0 GOOS=$(OS_VERSION) go build -v -a --ldflags="-s" -o ./out/bot-$(OS_VERSION) ./cmd/bot

run:
	$(DOCKER_IMG) /bin/sh -c '\
	go run ./cmd/bot/main.go \
	'

build:
	$(DOCKER_IMG) /bin/sh -c '\
	$(GO_GET_CMD) && \
	$(GO_TEST_CMD) && \
	$(GO_BUILD_CMD); \
	$(FIX_PERM) \
	'

build/win:
	$(DOCKER_IMG) /bin/sh -c '\
	$(GO_GET_CMD) && \
	$(GO_TEST_CMD) && \
	$(GO_BUILD_WIN_CMD) \
	'