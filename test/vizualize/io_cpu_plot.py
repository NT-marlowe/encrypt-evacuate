import json
import matplotlib.pyplot as plt
import sys


def plot(iostat_data, label: str):
    # データからタイムスタンプとCPU使用率を抽出
    xticks = [i for i in range(1, 15 + 1, 1)]
    cpu_usage = []

    for stat in iostat_data["sysstat"]["hosts"][0]["statistics"]:
        # timestamps.append(stat["timestamp"])
        avg_cpu = stat["avg-cpu"]
        cpu_usage.append(avg_cpu["user"] + avg_cpu["system"])

    # 折れ線グラフを作成
    plt.plot(
        xticks,
        cpu_usage,
        marker="o",
        label=label,
    )
    plt.xticks(xticks)
    # plt.tight_layout()


def visualize_cpu_usage():
    data_dir = "../result/io_cpu"
    files = ["baseline.json", "encryption_load.json", "proposed_method_load.json"]
    # JSONファイルを読み込む

    plt.figure(figsize=(12, 6))
    for file in files:
        path = f"{data_dir}/{file}"
        with open(path, "r") as f:
            data = json.load(f)

        # データからタイムスタンプとCPU使用率を抽出
        plot(data, file.split(".")[0])

    plt.xlabel("Timestamp")
    plt.ylabel("CPU Usage (user + system) [%]")
    plt.title("CPU Usage Over Time (User + System)")
    plt.grid(axis="y")
    plt.legend()
    plt.tight_layout()

    # グラフを表示
    plt.savefig("./img/cpu_usage.png")


# 実行例
visualize_cpu_usage()
