SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c


build:
	go build ./cmd/web

run:
	go run ./cmd/web

certs:
	@mkdir -p ./assets/tls/
	mkcert -key-file ./assets/tls/key.pem -cert-file ./assets/tls/cert.pem localhost 127.0.0.1 '::1'
