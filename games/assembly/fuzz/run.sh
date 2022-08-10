#!/bin/bash

echo "building..."
go-fuzz-build github.com/lmas/asm_game/fuzz

echo "running fuzzer..."
go-fuzz -bin=./fuzz-fuzz.zip -workdir=work
