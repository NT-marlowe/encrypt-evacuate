import json
import matplotlib.pyplot as plt


def parse_data(filename: str, metric: str):
    with open(filename, "r") as f:
        iostat_data = json.load(f)
        xticks = [
            i
            for i in range(1, len(iostat_data["sysstat"]["hosts"][0]["statistics"]) + 1)
        ]
        values = []

        for stat in iostat_data["sysstat"]["hosts"][0]["statistics"]:
            avg_cpu = stat["avg-cpu"]
            if metric == "user+system":
                values.append(float(avg_cpu["user"] + avg_cpu["system"]))
            else:
                values.append(float(avg_cpu[metric]))

    return xticks, values


def generate_filename_with_index(base_filename, index):
    """
    Generate a new filename by appending an index before the file extension.

    Args:
        base_filename (str): The base filename (e.g., 'cpu_usage.png').
        index (int): The index to append to the filename.

    Returns:
        str: New filename with the index appended (e.g., 'cpu_usage_1.png').
    """
    if ".png" not in base_filename:
        raise ValueError("The base filename must end with '.png'")

    # Split the filename into name and extension
    name, extension = base_filename.rsplit(".", 1)

    # Generate the new filename
    new_filename = f"{name}_{index}.{extension}"
    return new_filename


def visualize_cpu_metric(metric: str, output_file: str):
    """
    指定されたCPUメトリクスをプロットするグラフを作成し、保存
    :param metric: プロットするCPUメトリクス (例: 'user', 'system', 'iowait', 'user+system')
    :param output_file: グラフ画像の保存先ファイル名
    """
    data_dir = "../result/io_cpu"
    # data_dir = ".."
    files = ["baseline.json", "encryption_load.json", "proposed_method_load.json"]
    # files = ["encryption_load.json", "proposed_method_load.json"]
    iter = 5

    for i in range(1, iter + 1):
        plt.figure(figsize=(12, 6))
        for file in files:
            base_path = f"{data_dir}/{file}.{i}"
            # values, xticks = calc_average(base_path, metric, iter)
            xticks, values = parse_data(base_path, metric)
            if file == "proposed_method_load.json":
                print(values)

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
        plt.savefig(generate_filename_with_index(output_file, i))


# 実行例
# user + system の使用率をプロット
visualize_cpu_metric("user+system", "./img/cpu_usage.png")

# iowait をプロット
visualize_cpu_metric("iowait", "./img/iowait_usage.png")
