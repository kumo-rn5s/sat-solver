PROGRAM:=sat-solver

.PHONY: build
build:
	go build

.PHONY: ut
ut:
	go test

.PHONY: it
it: build
	time ./$(PROGRAM) test/sat/* | uniq
	time ./$(PROGRAM) test/unsat/* | uniq
