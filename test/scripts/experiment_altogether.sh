#!/bin/bash
#
# This script run the capturer once, and the ransomware program over all files.

set -e
set -u

# if not root user, exit
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root"
    exit
fi

# 実験するuser (whoami) に変更する
# つまり，k-akashiに変更する
# export USER=$(whoami) とするべきだったかも
USER=marlowe
EBPF_PROGRAM=ebpf-ssl

# パラメータ．デコードの並列処理を担当するスレッド (goroutine) の数を指定する．
parallelism=$1

cd $(git rev-parse --show-toplevel)
# 実行ファイルをビルドする．
make all
# 再現性のために`/data_shelter`内のファイルを削除する．ディレクトリ自体は削除しない．
make clear-shelter
cd test

# Fugaを実行する．
../${EBPF_PROGRAM} ${parallelism} &
pid=$!
sleep 1

# `openssl rand`コマンドで生成したランダムな内容のファイルが，test/data に用意されている．
# それらのファイルのパスをgrepで選定し，一個ずつmy_simple_ransomwareに渡す．
# .encファイルは暗号化後のファイルなので，除外する．
# 以下は1K, 10K, ..., 100MBのファイルを対象にしている例．
for file in $(ls ./data/power_10/* | grep -v enc); do
    # for file in $(ls ./data/power_2/* | grep -v enc); do
    # for file in $(ls ./data/incremental/10_1MB.data | grep -v enc); do
    # for file in $(ls ./data/2* | grep -v enc); do
    # file=$(basename $file)

    # my_simple_ransomwareは対象のファイルそのものは削除せず，暗号化したファイル (.enc) を新たに作成する．
    ./my_simple_ransomware ${file}
    echo "Ransomware ran on ${file}"

    sleep 4

    echo -----------

done

# FugaをSIGINTで終了させる．SIGKILLだと正常にclean upされないかもしれない．
kill -SIGINT $pid >/dev/null 2>&1
echo "Killed $pid"

# cmd/capturer のgoソースファイルの一部がroot:rootになるので，戻す．
# これをしないと，手動でmake allを叩いた時にpermission deniedになる．
# このスクリプトのみを使う場合は必ずしも必要ではない．
chown -R ${USER}:${USER} ../cmd/*
chown -R ${USER}:${USER} ../${EBPF_PROGRAM}
chown -R ${USER}:${USER} ./data/*

# rm ./data/*.enc

echo "parallelism: ${parallelism}"
