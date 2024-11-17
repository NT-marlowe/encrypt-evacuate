#!/usr/bin/env bash
#
set -eu
out_filename=$1
time_len=30

sudo sync
echo 3 | sudo tee /proc/sys/vm/drop_caches

iostat -xyt -o JSON 1 ${time_len} > ${out_filename} &

for file in $(ls ./files_sysbench | grep -v 'enc'); do
    echo "Encrypting ${file}"
    ./my_simple_ransomware ./files_sysbench/${file}
done

echo "Cleaning up..."
rm -rf ./files_sysbench/*.enc
