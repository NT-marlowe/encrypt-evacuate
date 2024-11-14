import os
import glob
import pandas as pd
import matplotlib.pyplot as plt


def extract_file_size(filename):
    """ファイル名からサイズ部分を抽出しMBに変換"""
    size = filename.split("_")[1].split(".")[0]
    return int(size[:-2])


def load_and_merge_data(file_pattern):
    """複数ファイルからデータを読み込み、辞書形式で保持"""
    data = {}
    for filepath in glob.glob(file_pattern):
        key = os.path.basename(filepath).split("_")[1].split(".")[0]  # p[数字]を取得
        df = pd.read_csv(filepath, skiprows=1, header=None, names=["filename", "value"])
        df["size"] = df["filename"].apply(extract_file_size)
        data[key] = df
    return data


def plot_combined_data(data, title):
    """複数のデータを1つのグラフにまとめてプロット"""
    plt.figure(figsize=(10, 6))
    data = dict(sorted(data.items(), key=lambda x: int(x[0])))
    for key, df in data.items():
        plt.plot(df["size"], df["value"], marker="o", label=key)
        print(key)
        plt.legend()

    plt.xlabel("File Size (MB)")
    plt.xticks([i for i in range(1, 11, 1)])
    plt.ylabel("Value")
    plt.title(f"{title} Plot")
    # plt.legend(title="p[数字]")
    plt.grid(True)
    # plt.show()
    plt.savefig(f"./img/{title.lower()}.png")


# reten ファイルを1つのグラフにプロット
reten_data = load_and_merge_data("../result/measure_parallelism_seek/reten_*")
plot_combined_data(reten_data, "Reten_Seek")

# match ファイルを1つのグラフにプロット
match_data = load_and_merge_data("../result/measure_parallelism_seek/match_*")
plot_combined_data(match_data, "Match_Seek")
