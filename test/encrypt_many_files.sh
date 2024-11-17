#!/usr/bin/env bash
#

for file in $(ls ./files_sysbench | grep -v 'enc'); do
    echo "Encrypting ${file}"
    ./my_simple_ransomware ./files_sysbench/${file}
done

echo "Cleaning up..."
rm -rf ./files_sysbench/*.enc
