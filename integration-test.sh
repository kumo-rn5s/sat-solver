#!/bin/sh -eu

readonly FILE="integration-test"
time ./sat-solver test/sat/* | uniq | tee "$FILE"
time ./sat-solver test/unsat/* | uniq | tee -a "$FILE"

nlines=$(wc -l "$FILE" | awk '{print $1}')
if [ "$nlines" = 2 ]; then
  echo "Integration Test Successfully"
else
  echo "Integration Test Failure"
  exit 1
fi
