#!/bin/bash
#
# Helper script for working with fuzzit.dev:
#
# https://github.com/fuzzitdev/example-go
#
# Based on:
#
# https://github.com/google/syzkaller/blob/master/fuzzit.sh
set -eux

function target {
    go-fuzz-build -libfuzzer -func $3 -o fuzzer.a $2
    clang -fsanitize=fuzzer fuzzer.a -o fuzzer

    ./fuzzit create target $1 || true
    ./fuzzit create job $LOCAL --type fuzzing --branch $TRAVIS_BRANCH --revision $TRAVIS_COMMIT prose/$1 ./fuzzer
}

go get -u github.com/dvyukov/go-fuzz/go-fuzz-build
wget -q -O fuzzit https://github.com/fuzzitdev/fuzzit/releases/download/v2.4.12/fuzzit_Linux_x86_64
chmod a+x fuzzit

./fuzzit auth $FUZZIT_API_KEY
if [ "$1" = "fuzzing" ]; then
    export LOCAL=""
else
    export LOCAL="--local"
fi

target prose-transform ./transform Fuzz
target prose-tokenize ./tokenize Fuzz
target prose-summarize ./summarize Fuzz
target prose-chunk ./chunk Fuzz
