#!/bin/bash

set -e
set -u


if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# parallelism=$1
parallelism=("1" "2" "4" "8" "12" "15" "16")
for p in ${parallelism[@]}; do
    ./experiment.sh ${p}
    ./calc_reten_match.sh reten > p${p}_reten.txt
    ./calc_reten_match.sh match > p${p}_match.txt

done
    
#     ./calc_reten_match.sh reten > p${parallelism}_${i}_reten.txt
#     ./calc_reten_match.sh match > p${parallelism}_${i}_match.txt
# done




# for i in $(seq 1 5); do
#     ./experiment.sh ${parallelism}
    
#     ./calc_reten_match.sh reten > p${parallelism}_${i}_reten.txt
#     ./calc_reten_match.sh match > p${parallelism}_${i}_match.txt
# done
