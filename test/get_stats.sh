#!/bin/bash

# if not root user, exit
if [ "$(id -u)" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

set -e
set -u
# cpu usafe, memory usage, max memory usage
# percent of cpu
# maximum RSS
# averate RSS
# file system Output

datafile=$1
statsfile=./result/stats_$(basename $datafile .data).txt

# echo "cpu usage, max RSS (kB), averate total memory, file system output" > ${statsfile}
# /usr/bin/time -f "%P, \t%M, \t%K, \t%O" -a -o ${statsfile} ../ebpf-ssl tmp &
/usr/bin/time -a -o ${statsfile} -v ../ebpf-ssl tmp &
time_pid=$!

ebpf_pid=$(pgrep -P $time_pid) 
ps auxww | grep ebpf | grep -v grep

sleep 1

./my_simple_ransomware ${datafile} -e


echo "killing ebpf: $ebpf_pid"
kill -SIGINT $ebpf_pid 

sleep 1

# if ps -p $ebpf_pid > /dev/null; then
    # kill -SIGKILL $ebpf_pid
# fi
ps auxww | grep ebpf | grep -v grep

wait $time_pid
