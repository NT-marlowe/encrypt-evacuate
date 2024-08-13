#!/bin/bash
#
set -e


# echo "filename, levenstein_distance, partial_match_ratio"
echo "filename, match_ratio"
# for file in $(ls ./data | grep -v enc); do
for file in $(ls ./data/1* | grep -v enc); do
    file=$(basename $file)
    echo -n "$file, "

    original=./data/$file
    recovered=/usr/tmp/data_shelter/$file
    python3 commands.py match ${original} ${recovered}
done


