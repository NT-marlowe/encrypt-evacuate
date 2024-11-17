import json
import matplotlib.pyplot as plt


def plot(iostat_data, label: str, metric: str):
    """
    指定されたCPUメトリクスをプロットする
    :param iostat_data: JSON形式のiostatデータ
    :param label: グラフの凡例に使用するラベル
    :param metric: プロットするCPUメトリクス (例: 'user', 'system', 'iowait')
    """
    xticks = [
        i for i in range(1, len(iostat_data["sysstat"]["hosts"][0]["statistics"]) + 1)
    ]
    values = []

    for stat in iostat_data["sysstat"]["hosts"][0]["statistics"]:
        avg_cpu = stat["avg-cpu"]
        if metric == "user+system":
            values.append(avg_cpu["user"] + avg_cpu["system"])
        else:
            values.append(avg_cpu.get(metric, 0))

    # 折れ線グラフを作成
    plt.plot(
        xticks,
        values,
        marker="o",
        label=label,
    )
    plt.xticks(xticks)


def visualize_cpu_metric(metric: str, output_file: str):
    """
    指定されたCPUメトリクスをプロットするグラフを作成し、保存
    :param metric: プロットするCPUメトリクス (例: 'user', 'system', 'iowait', 'user+system')
    :param output_file: グラフ画像の保存先ファイル名
    """
    data_dir = "../result/io_cpu"
    files = ["baseline.json", "encryption_load.json", "proposed_method_load.json"]

    plt.figure(figsize=(12, 6))
    for file in files:
        path = f"{data_dir}/{file}"
        with open(path, "r") as f:
            data = json.load(f)

        # プロット
        plot(data, file.split(".")[0], metric)

    # ラベル設定
    plt.xlabel("Timestamp")
    plt.ylabel(
        f"{metric.capitalize()} [%]"
        if metric != "user+system"
        else "CPU Usage (user + system) [%]"
    )
    plt.title(f"{metric.capitalize()} Over Time")
    plt.grid(axis="y")
    plt.legend()
    plt.tight_layout()

    # グラフを保存
    plt.savefig(output_file)


# 実行例
# user + system の使用率をプロット
visualize_cpu_metric("user+system", "./img/cpu_usage.png")

# iowait をプロット
visualize_cpu_metric("iowait", "./img/iowait_usage.png")
