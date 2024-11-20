#!/usr/bin/env bash
#
set -eu
out_filename=$1

for i in $(seq 1 5); do
    echo "Running iteration $i"
    echo "Removing caches..."
    sudo sync
    echo 3 | sudo tee /proc/sys/vm/drop_caches

    iostat -xyt -o JSON 1 30 >${out_filename}.${i}
    # pid=$!

    # wait $pid
    echo "Iteration ${i} done. Sleeping for 5 seconds..."
    sleep 5
done
