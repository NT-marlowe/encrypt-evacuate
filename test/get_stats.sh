#!/bin/bash

# if not root user, exit
if [ "$(id -u)" -ne 0 ]; then
    echo "Please run as root"
    exit 1
fi

set -e
# cpu usafe, memory usage, max memory usage
# percent of cpu
# maximum RSS
# averate RSS
# file system Output

echo "cpu usage, max RSS (kB), averate total memory, file system output" > ./result/stats.txt
# /usr/bin/time -f "%P, \t%M, \t%K, \t%O" -o ./result/stats.txt ../ebpf-ssl tmp &
/usr/bin/time -a -o ./result/stats.txt -v ../ebpf-ssl tmp &
time_pid=$!

ebpf_pid=$(pgrep -P $time_pid) 
echo $ebpf_pid

./my_simple_ransomware ./data/04_10MB.data -e

kill -SIGINT $ebpf_pid 

if ps -p $ebpf_pid > /dev/null; then
    kill -SIGKILL $ebpf_pid
fi

wait $time_pid
