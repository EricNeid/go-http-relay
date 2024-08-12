<!--
SPDX-FileCopyrightText: 2021 Eric Neidhardt
SPDX-License-Identifier: CC-BY-4.0
-->
<!-- markdownlint-disable MD041-->
[![Go Report Card](https://goreportcard.com/badge/github.com/EricNeid/go-webserver?style=flat-square)](https://goreportcard.com/report/github.com/EricNeid/go-webserver)
![Go](https://github.com/EricNeid/go-webserver/workflows/Go/badge.svg)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/EricNeid/go-webserver)
[![Release](https://img.shields.io/github/release/EricNeid/go-webserver.svg?style=flat-square)](https://github.com/EricNeid/go-webserver/releases/latest)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/EricNeid/go-webserver)

# About

A super simple webserver, written in go.

Based on: <https://marcofranssen.nl/go-webserver-with-graceful-shutdown/>

## Quickstart

The simples way to compile this application is to use the provide makefile.
It provides cross compilation to linux and windows and makes use of docker.

Docker:

```bash
make build-windows
make build-linux
```

Manual and without docker:

```bash
go build -o ./out/ ./cmd/mapprovider/
```

Start server:

```bash
./webserver -listen-addr=:80 -base-path=foo
```

## Options

Application can be configure using command line arguments or
environment variables or a combination of both.

* listen-addr/LISTEN_ADDR - listing address, ie. ":5000"
* base-path/BASE_PATH - base path to serve application, ie "/custom"

Example:

```bash
./webserver -base-path webserver-0.1.0 -listen-addr :8080
```

## Endpoints

The following endpoints are available, assuming
the configuration is not changed.

* <http://localhost:5000>

## Testing

To run tests:

```bash
make test
```

## Question or comments

Please feel free to open a new issue:
<https://github.com/EricNeid/go-webserver/issues>
