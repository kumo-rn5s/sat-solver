.PHONY: test
test:
	go run main.go test/sat/* | uniq
	go run main.go test/unsat/* | uniq