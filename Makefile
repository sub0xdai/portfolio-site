.PHONY: build run templ

TEMPL_BIN=/home/sub0x/.asdf/installs/golang/1.23.2/packages/bin/templ

templ:
	$(TEMPL_BIN) generate

build: templ
	go build -o bin/resume-ai cmd/main.go

run: build
	./bin/resume-ai

dev:
	$(TEMPL_BIN) generate --watch & \
	go run cmd/main.go
