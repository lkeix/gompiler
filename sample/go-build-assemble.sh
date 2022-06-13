#!/bin/bash
from=$1
to=$2

go tool compile -S -N $from > $to