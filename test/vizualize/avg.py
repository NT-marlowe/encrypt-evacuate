import glob
import pandas as pd
import matplotlib.pyplot as plt


def extract_file_size(filename):
    """ファイル名からサイズ部分を抽出しMBに変換"""
    size = filename.split("_")[1].split(".")[0]
    return int(size[:-2])


def load_data_and_average(file_pattern):
    """複数のファイルからデータを読み込み、同じ行の数値を平均化"""
    all_data = []

    # 指定されたパターンのファイルをすべて読み込む
    n = 0
    for filepath in glob.glob(file_pattern):
        df = pd.read_csv(filepath, skiprows=1, header=None, names=["filename", "value"])
        df["size"] = df["filename"].apply(extract_file_size)
        all_data.append(df[["size", "value"]])
        n += 1

    # サイズごとに数値の平均を計算
    merged_data = pd.concat(all_data)
    averaged_data = merged_data.groupby("size").mean().reset_index()

    return averaged_data, n


def plot_average_data(reten_data, match_data, num_reten_file, num_match_file):
    """retenとmatchのデータをプロット"""
    plt.figure(figsize=(10, 6))

    # retenデータのプロット
    plt.plot(
        reten_data["size"],
        reten_data["value"],
        marker="o",
        label=f"Reten Avg. ({num_reten_file} files)",
        color="blue",
    )

    # matchデータのプロット
    plt.plot(
        match_data["size"],
        match_data["value"],
        marker="o",
        label=f"Match Avg. ({num_match_file} files)",
        color="red",
    )

    plt.xlabel("File Size (MB)")
    plt.xticks([i for i in range(1, 11, 1)])
    plt.ylabel("Average Value")
    plt.title("Average Reten and Match Values")
    plt.legend()
    plt.grid(True)
    # plt.show()
    plt.savefig("./img/avg.png")


# ファイルパスパターンに従ってデータを読み込み、平均を計算
reten_data, num_reten_file = load_data_and_average(
    "../result/measure_parallelism/p15*_reten.txt"
)
match_data, num_match_file = load_data_and_average(
    "../result/measure_parallelism/p15*_match.txt"
)

# プロット
plot_average_data(reten_data, match_data, num_reten_file, num_match_file)
