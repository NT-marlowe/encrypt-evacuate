#!/usr/bin/env bash

set -ue

function main() {

    local current_dir
    current_dir=$(dirname "${BASH_SOURCE[0]:-$0}")
    source ${current_dir}/util.sh
    local result_dir=${current_dir}/../result/measure_parallelism_seek

    check_root

    parallelisms=(1 2 4 6 8 10 12 14 16 24 32)
    # parallelisms=(1)
    for p in ${parallelisms[@]}; do

        ${current_dir}/experiment_altogether.sh ${p}

        padded_p=$(printf "%02d" ${p})

        ${current_dir}/../calc_reten_match.sh reten >${result_dir}/reten_${padded_p}.csv
        ${current_dir}/../calc_reten_match.sh match >${result_dir}/match_${padded_p}.csv

        sleep 2
    done

    set_owner ${result_dir}
}

main "$@"
