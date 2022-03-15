#!/bin/sh
set -e

RESULT_SAT=$(./sat-solver test/sat/* | uniq)
RESULT_UNSAT=$(./sat-solver test/unsat/* | uniq)

if [ \( "$RESULT_SAT" != "sat" \) ] || [ \( "$RESULT_UNSAT" != "unsat" \) ]; then
    exit 1
else
    echo "Integration Test Successfully"
fi
