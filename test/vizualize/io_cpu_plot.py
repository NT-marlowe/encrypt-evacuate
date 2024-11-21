import json
import matplotlib.pyplot as plt


def parse_data(iostat_data, metric: str):
    xticks = [
        i for i in range(1, len(iostat_data["sysstat"]["hosts"][0]["statistics"]) + 1)
    ]
    values = []

    for stat in iostat_data["sysstat"]["hosts"][0]["statistics"]:
        avg_cpu = stat["avg-cpu"]
        if metric == "user+system":
            values.append(avg_cpu["user"] + avg_cpu["system"])
        else:
            values.append(avg_cpu[metric])

    return xticks, values


def calc_average(base_filename: str, metric: str, iter: int):
    """
    指定されたmetricsの値の配列を取得し，iteration間の平均値を計算する．
    """
    val_2d_array = []
    for i in range(1, iter + 1):
        path = f"{base_filename}.{i}"
        with open(path, "r") as f:
            data = json.load(f)
            xticks, values = parse_data(data, metric)
            val_2d_array.append(values)

    res = []
    series_len = len(xticks)
    for j in range(series_len):
        tmp = 0
        for i in range(iter):
            tmp += val_2d_array[i][j]

        res.append(tmp / series_len)

    return res, xticks


def visualize_cpu_metric(metric: str, output_file: str):
    """
    指定されたCPUメトリクスをプロットするグラフを作成し、保存
    :param metric: プロットするCPUメトリクス (例: 'user', 'system', 'iowait', 'user+system')
    :param output_file: グラフ画像の保存先ファイル名
    """
    data_dir = "../result/io_cpu"
    files = ["baseline.json", "encryption_load.json", "proposed_method_load.json"]
    iter = 5

    plt.figure(figsize=(12, 6))
    for file in files:
        base_path = f"{data_dir}/{file}"
        values, xticks = calc_average(base_path, metric, iter)

        plt.plot(
            xticks,
            values,
            marker="o",
            label=file.split(".")[0],
        )
        plt.xticks(xticks, xticks)

    # ラベル設定
    plt.xlabel("Timestamp")
    plt.ylabel(
        f"{metric.capitalize()} [%]"
        if metric != "user+system"
        else "CPU Usage (user + system) [%]"
    )
    plt.title(f"{metric.capitalize()} Over Time ({iter} times average)")
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
