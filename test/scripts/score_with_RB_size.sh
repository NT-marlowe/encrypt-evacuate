#!/usr/bin/env bash

set -ue

function main() {
    local current_dir
    current_dir=$(dirname "${BASH_SOURCE[0]:-$0}")
    source ${current_dir}/util.sh

    check_root

    local parallelism=$1
    local RB_size=$2

    # ${current_dir}/experiment_altogether.sh ${parallelism}

    padded_p=$(printf "%02d" ${parallelism})

    local result_dir=${current_dir}/../result/buf_capability/journal
    # ${current_dir}/../calc_reten_match.sh reten >${result_dir}/reten_${padded_p}.csv
    ${current_dir}/../calc_reten_match.sh match >${result_dir}/match_${padded_p}_${RB_size}MiB.csv

    sleep 2

    set_owner ${result_dir}
}

main "$@"
