SHELL := /bin/sh

GOOS ?= linux
GOARCH ?= amd64

.PHONY: all clean

all: buttstickerapi

clean:
	$(RM) buttstickerapi

buttstickerapi: cmd/buttstickerapi/buttstickerapi.go internal/pkg/handler/TickerHandler.go
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $@ $<
