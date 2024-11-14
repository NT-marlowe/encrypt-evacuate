#!/usr/bin/env bash

set -ue

USER=marlowe

function check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo "Please run as root"
        exit
    fi
}

function set_owner() {
    local path=$1

    chown -R ${USER}:${USER} ${path}
}
