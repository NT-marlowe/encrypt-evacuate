import json
import matplotlib.pyplot as plt
import sys


def visualize_cpu_usage(json_file):
    # JSONファイルを読み込む
    with open(json_file, "r") as file:
        data = json.load(file)

    # データからタイムスタンプとCPU使用率を抽出
    timestamps = []
    cpu_usage = []

    for stat in data["sysstat"]["hosts"][0]["statistics"]:
        vda_disk_data = stat["disk"][-1]
        # print(stat["disk"][-1]["wkB/s"], end=", ")
        print(f"rkB/s = {vda_disk_data['rkB/s']}\t\twkB/s = {vda_disk_data['wkB/s']}")
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
    json_file = sys.argv[1]
    visualize_cpu_usage(json_file)
