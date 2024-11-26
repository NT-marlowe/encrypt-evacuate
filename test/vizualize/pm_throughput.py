import json
import sys
import numpy as np


def load_disk_data(json_file):
    # JSONファイルを読み込む
    with open(json_file, "r") as file:
        data = json.load(file)

    res = []

    # write_data = []
    for stat in data["sysstat"]["hosts"][0]["statistics"]:
        sda_disk_data = stat["disk"][-1]
        # print(stat["disk"][-1]["wkB/s"], end=", ")
        rkB_s = sda_disk_data["rkB/s"]
        wkB_s = sda_disk_data["wkB/s"]
        disk_util = sda_disk_data["util"]
        res.append((rkB_s, wkB_s, disk_util))

    #     if wkB_s > 0:
    #         write_data.append(wkB_s)

    # print(sorted(write_data))
    return res


def print_diff(json_file1, json_file2):
    print(f"{json_file2} - {json_file1}")

    disk_data_1 = load_disk_data(json_file1)
    disk_data_2 = load_disk_data(json_file2)
    for i in range(len(disk_data_1)):
        rkB_s_diff = disk_data_2[i][0] - disk_data_1[i][0]
        wkB_s_diff = disk_data_2[i][1] - disk_data_1[i][1]
        util_diff = disk_data_2[i][2] - disk_data_1[i][2]
        print(
            f"rkB/s = {rkB_s_diff:>10}\t wkB/s = {wkB_s_diff:>10}\t util = {util_diff:>6.2f}"
        )


def accumulate_positive_write(json_file1, json_file2):
    print(f"{json_file2} - {json_file1}")

    write_data_MBs = []
    for idx in range(1, 6, 1):
        disk_data_1 = load_disk_data(f"{json_file1}.{idx}")
        disk_data_2 = load_disk_data(f"{json_file2}.{idx}")
        for i in range(len(disk_data_1)):
            # rkB_s_diff = disk_data_2[i][0] - disk_data_1[i][0]
            wkB_s_diff = disk_data_2[i][1] - disk_data_1[i][1]
            # util_diff = disk_data_2[i][2] - disk_data_1[i][2]
            # if wkB_s_diff > 1000:
            if abs(wkB_s_diff) > 500:  # more than 1MB/s
                # if abs(wkB_s_diff) > 100:
                write_data_MBs.append(wkB_s_diff / 1000)

    return write_data_MBs


#     return write_data_MBs
def print_stat(write_data_MBs):
    print((write_data_MBs))
    print(sorted(write_data_MBs))
    print(f"Mean: {np.mean(write_data_MBs)} [MB/s]")
    print(f"Median: {np.median(write_data_MBs)} [MB/s]")
    print(f"Max: {np.max(write_data_MBs)} [MB/s]")
    print(f"Min: {np.min(write_data_MBs)} [MB/s]")
    print(f"Std: {np.std(write_data_MBs)} [MB/s]")


if __name__ == "__main__":
    if len(sys.argv) == 2:
        for i in range(1, 6):
            load_disk_data(f"{sys.argv[1]}.{i}")
        exit(0)
    elif len(sys.argv) == 3:
        json_file1, json_file2 = sys.argv[1], sys.argv[2]
        write_data_MBs = accumulate_positive_write(json_file1, json_file2)
        print_stat(write_data_MBs)
