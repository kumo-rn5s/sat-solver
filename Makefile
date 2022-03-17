CNF:=test/sat/uf100-01.cnf
IMAGE:=sat-solver
PROGRAM:=sat-solver
TAG:=$(shell git rev-parse HEAD)

.PHONY: build
build:
	go build "-ldflags=-s -w -buildid=" -trimpath -o $(PROGRAM)

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE):$(TAG) .

.PHONY: docker-run
docker-run: docker-build
	docker run --rm -i $(IMAGE):$(TAG) < $(CNF)

.PHONY: docker-run-test
docker-run-test: docker-unit-test docker-integration-test

.PHONY: docker-unit-test
docker-unit-test: docker-build
	docker run --rm $(IMAGE):$(TAG) go test

.PHONY: docker-integration-test
docker-integration-test: docker-build
	docker run --rm $(IMAGE):$(TAG) ./integration-test.sh

.PHONY: clean
clean:
	rm $(PROGRAM)
