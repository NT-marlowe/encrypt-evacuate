#!/bin/bash

set -e
set -u

# base_num=2
subcom=$1

power_2=(2 4 8 16 32)
for i in ${power_2[@]}; do
    echo "Running with ${i}MiB ring buffer"
    ./calc_match.sh ${subcom} $i > ./result/buf_capability/${subcom}_${i}MiB.csv
done
