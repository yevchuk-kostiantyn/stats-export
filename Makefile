export CGO_ENABLED=0

SVC_NAME = stats-export

BUILD_DIR=.build

.PHONY: build clean local-env

all: test build run

verify:
	go vet ./...

## test: tests go service
test: verify
	go test ./...

get-deps:
	env GIT_TERMINAL_PROMPT=1 go mod download

## build: build go service
build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(SVC_NAME) .

## run: runs service locally
run:
	env $$(cat $(SVC_NAME).env.local) $(BUILD_DIR)/$(SVC_NAME)

help: Makefile
	@echo " Choose a command run in \033[32m"$(SVC_NAME)"\033[0m:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

