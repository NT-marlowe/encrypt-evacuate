#!/usr/bin/env bash

set -ue

function main() {
    local current_dir
    current_dir=$(dirname "${BASH_SOURCE[0]:-$0}")
    source ${current_dir}/util.sh

    check_root

    local map_file=${current_dir}/../../ebpf_src/maps.h
    local result_dir=${current_dir}/../result/buf_capability/journal

    local parallelism=$1
    local padded_p=$(printf "%02d" ${parallelism})

    local RB_sizes=(2 4 8 16 32)
    for size in ${RB_sizes[@]}; do
        echo "RB size: ${size} MiB"
        sed -i "s/^#define MULT .*$/#define MULT $((size))/" ${map_file}

        ${current_dir}/experiment_altogether.sh ${parallelism}
        ${current_dir}/../calc_reten_match.sh match >${result_dir}/match_${padded_p}_${size}MiB.csv
    done

    # ${current_dir}/../calc_reten_match.sh reten >${result_dir}/reten_${padded_p}.csv

    sleep 2

    sed -i "s/^#define MULT .*$/#define MULT 1/" ${map_file}
    set_owner ${result_dir}
}

main "$@"
