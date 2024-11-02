#!/bin/bash

set -e
set -u

iter=10000
filename=$1
function iterate_encrypt() {
    for i in $(seq 1 ${iter}); do
        ./foobar ./data/03_1MB.data
    done
}

{ time iterate_encrypt; } > $filename 2>&1
