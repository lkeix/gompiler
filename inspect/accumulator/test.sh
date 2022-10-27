#!/bin/bash


assert() {
  expect="$1"
  input="$2"
  
  go run main.go -value=$2 > main.s
  cc -o main main.s
  ./main -value=1
  
  actual="$?"

  if [ "$actual" = "$expect" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expect expect, but got $actual"
    exit 1
  fi
}

assert 0 0
assert 42 42

echo ok

rm -rf main.s main
