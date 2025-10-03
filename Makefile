.PHONY: test build

VERSION = $(shell grep "const Version =" cmd/schemaguesser/cmd/version.go | grep "const Version =" | sed -e 's-.*= `--' -e 's-`--')
BIN_UPLOAD_PATH =
SHELL=/bin/bash

build:
	go build -o build/mschemaguesser -ldflags "-s -w" cmd/schemaguesser/main.go

build-mcopy-docker:
	bash -c ./buildImage.sh

tag-release:
	git tag -a v$(VERSION) -m "Release $(VERSION)"
	git push --tags

start-dev-env:
	bash -c "docker/bin/compose_env.sh start"

stop-dev-env:
	bash -c "docker/bin/compose_env.sh stop"

restart-dev-env:
	bash -c "docker/bin/compose_env.sh stop"
	bash -c "docker/bin/compose_env.sh start"

test:
	go test -cover ./... && echo ":)" || echo ":-/"

test-without-it:
	go test --skip "_IT" ./... && echo ":)" || echo ":-/"


test-integration:
	docker/bin/test_env.sh test
