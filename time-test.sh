#!/bin/sh
set -e

time ./sat-solver test/sat/* | uniq
time ./sat-solver test/unsat/* | uniq
