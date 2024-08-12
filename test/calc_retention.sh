#!/bin/bash

set -e


echo "filename, retention_rate"

for file in $(ls ./data | grep -v enc); do
    echo -n "$file, "
    original=./data/$file
    recovered=/usr/tmp/data_shelter/$file
    
    python3 commands.py reten ${original} ${recovered}
done
