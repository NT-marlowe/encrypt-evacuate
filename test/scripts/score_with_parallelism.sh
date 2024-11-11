#!/usr/bin/env bash

set -ue

function main() {
    local current_dir
    current_dir=$(dirname "${BASH_SOURCE[0]:-$0}")
    source ${current_dir}/util.sh

    check_root

    echo ${current_dir}

    parallelisms=(1 2)
    for p in ${parallelisms[@]}; do
        ${current_dir}/../experiment_altogether.sh ${p}
    done

}

main "$@"
