#!/usr/bin/env bash

set -ue

function main() {
    local current_dir
    current_dir=$(dirname "${BASH_SOURCE[0]:-$0}")
    source ${current_dir}/util.sh

    check_root

    parallelisms=(1 2 4 6 8 10 15 16 32)
    for p in ${parallelisms[@]}; do

        ${current_dir}/experiment_altogether.sh ${p}

        padded_p=$(printf "%02d" ${p})

        echo
        pwd
        echo

        ${current_dir}/../calc_reten_match.sh reten >${current_dir}/../result/reten_${padded_p}.csv
        ${current_dir}/../calc_reten_match.sh match >${current_dir}/../result/match_${padded_p}.csv
    done
}

main "$@"
