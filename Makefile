.PHONY: build buid-tests unit-tests test

.DEFAULT_GOAL: build

build:
	GO111MODULE=on go mod tidy
	docker build -t krane .

build-tests:
	docker build -t tests/krane --target build .

unit-tests: build-tests
	docker run --rm -it tests/krane go test -v ./...

test: build
	docker run --rm -it -v ${PWD}/tests:/krane/tests -e LOG_LEVEL=DEBUG -e CI=true -e KRANE_ROOT=/krane krane 