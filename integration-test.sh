#!/bin/sh -eu

test() {
  if [ "$1" != "$2" ]; then
    echo "Integration Test Failure"
    exit 1
  fi
}

main() {
  result=$(./sat-solver test/sat/* | uniq)
  test sat "$result"
  result=$(./sat-solver test/unsat/* | uniq)
  test unsat "$result"
  echo "Integration Test Successfully"
}

main
