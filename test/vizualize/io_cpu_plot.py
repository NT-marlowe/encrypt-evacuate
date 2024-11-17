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
    plt.figure(figsize=(12, 6))
    plt.plot(
        xticks,
        cpu_usage,
        marker="o",
        label=label,
    )
    plt.xticks()
    plt.tight_layout()


def visualize_cpu_usage(json_file):
    # JSONファイルを読み込む
    with open(json_file, "r") as file:
        data = json.load(file)

    # データからタイムスタンプとCPU使用率を抽出
    plot(data, "CPU Usage")

    plt.xticks(rotation=45)
    plt.xlabel("Timestamp")
    plt.ylabel("CPU Usage (%)")
    plt.title("CPU Usage Over Time (User + System)")
    plt.grid(True)
    plt.legend()
    plt.tight_layout()

    # グラフを表示
    plt.savefig("./img/cpu_usage.png")


# 実行例
visualize_cpu_usage(sys.argv[1])
