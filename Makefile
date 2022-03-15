IMAGE:=sat-solver
TAG:=latest

.PHONY: build
build:
	go build

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE):$(TAG) .

.PHONY: docker-run-test
docker-run-test: docker-unit-test docker-integration-test

.PHONY: docker-unit-test
docker-unit-test:
	docker run --rm -v $(PWD):/src --workdir=/src golang:1.17.8-alpine3.15 sh -c "go test"

.PHONY: docker-integration-test
docker-integration-test:
	docker run --rm --entrypoint="" -v $(PWD)/test:/test $(IMAGE):$(TAG) ./integration-test.sh

.PHONY: docker-timer-test
docker-timer-test: 
	docker run --rm --entrypoint="" -v $(PWD)/test:/test $(IMAGE):$(TAG) ./time-test.sh
