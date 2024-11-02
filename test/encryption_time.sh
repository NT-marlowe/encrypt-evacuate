#!/bin/bash

set -e
set -u

iter=$1
filename=$2
function iterate_encrypt() {
    for i in $(seq 1 ${iter}); do
        ./foobar ./data/03_1MB.data
    done
}

{ time iterate_encrypt; } > $filename 2>&1
