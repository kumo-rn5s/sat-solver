PROGRAM:=sat-solver

.PHONY: build
build:
	go build

.PHONY: test
test: build
	time ./$(PROGRAM) test/sat/* | uniq
	time ./$(PROGRAM) test/unsat/* | uniq
