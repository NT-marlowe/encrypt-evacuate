#!/usr/bin/env bash
#
set -eu
out_filename=$1

iostat -xyt -o JSON 1 15 > ${out_filename} &

for file in $(ls ./files_sysbench | grep -v 'enc'); do
    echo "Encrypting ${file}"
    ./my_simple_ransomware ./files_sysbench/${file}
done

echo "Cleaning up..."
rm -rf ./files_sysbench/*.enc
