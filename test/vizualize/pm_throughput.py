import json
import matplotlib.pyplot as plt
import sys


def load_disk_data(json_file):
    # JSONファイルを読み込む
    with open(json_file, "r") as file:
        data = json.load(file)

    res = []

    for stat in data["sysstat"]["hosts"][0]["statistics"]:
        vda_disk_data = stat["disk"][-1]
        # print(stat["disk"][-1]["wkB/s"], end=", ")
        rkB_s = vda_disk_data["rkB/s"]
        wkB_s = vda_disk_data["wkB/s"]
        disk_util = vda_disk_data["util"]
        # print(
        #     # f"rkB/s = {vda_disk_data['rkB/s']:>10}\t wkB/s = {vda_disk_data['wkB/s']:>10}"
        #     f"rkB/s = {rkB_s:>10}\t wkB/s = {wkB_s:>10}\t util = {disk_util:>10}"
        # )
        res.append((rkB_s, wkB_s, disk_util))

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
    # exit(0)
    #     timestamps.append(stat["timestamp"])
    #     avg_cpu = stat["avg-cpu"]
    #     cpu_usage.append(avg_cpu["user"] + avg_cpu["system"])

    # # 折れ線グラフを作成
    # plt.figure(figsize=(12, 6))
    # plt.plot(timestamps, cpu_usage, marker="o", label="CPU Usage (User + System)")
    # plt.xticks(rotation=45)
    # plt.xlabel("Timestamp")
    # plt.ylabel("CPU Usage (%)")
    # plt.title("CPU Usage Over Time (User + System)")
    # plt.grid(True)


if __name__ == "__main__":
    json_file1, json_file2 = sys.argv[1], sys.argv[2]
    print_diff(json_file1, json_file2)
