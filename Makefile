PROGRAM:=sat-solver

build:
	go build

.PHONY: test
test: build
	time ./$(PROGRAM) test/sat/* | uniq
	time ./$(PROGRAM) test/unsat/* | uniq
