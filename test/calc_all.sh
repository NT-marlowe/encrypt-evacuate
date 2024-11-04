#!/bin/bash

set -e
set -u

# base_num=2
subcom=$1

power_2=(2 4 8 16 32 64)
for p2 in ${power_2[@]}; do
    echo "Running with ${p2}MiB ring buffer"
    p2_padded=$(printf "%02d" $p2)
    ./calc_reten_match.sh ${subcom} ${p2_padded} | tee ./result/buf_capability/${p2_padded}MiB_${subcom}.csv
done
