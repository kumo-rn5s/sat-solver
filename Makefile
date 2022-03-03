PROGRAM:=sat-solver

.PHONY: build
build:
	go build


.PHONY: ut
ut: build
	go test -v ./...

.PHONY: it
it: build
	time ./$(PROGRAM) test/sat/* | uniq
	time ./$(PROGRAM) test/unsat/* | uniq
