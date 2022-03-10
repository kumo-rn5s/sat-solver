PROGRAM:=sat-solver

.PHONY: build
build:
	go build

.PHONY: test
test: unit-test integration-test

.PHONY: unit-test
unit-test:
	go test

.PHONY: integration-test
integration-test:
	time ./$(PROGRAM) test/sat/* | uniq
	time ./$(PROGRAM) test/unsat/* | uniq
