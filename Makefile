# SPDX-FileCopyrightText: 2021 Eric Neidhardt
# SPDX-License-Identifier: CC0-1.0

DIR := ${CURDIR}
GO_IMAGE := golang:1.22-alpine
LINTER_IMAGE := golangci/golangci-lint:v1.54-alpine

.PHONY: build-windows
build-windows:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \ \
		go build -o ./out/ ./cmd/webserver/
.PHONY: build-linux
build-linux:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \ \
		go build -o ./out/ ./cmd/webserver/


.PHONY: cover
cover:
	docker run -it --rm \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		mkdir -p out && go test -coverprofile=out/cover.out ./... && go tool cover -html=out/cover.out

.PHONY: test
test:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go test ./...

.PHONY: lint
lint:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		${LINTER_IMAGE} \
		golangci-lint --timeout=5m run ./...


.PHONY: clean
clean:
	rm -rf out

.PHONY: format
format:
	docker run -it --rm \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go fmt ./...
