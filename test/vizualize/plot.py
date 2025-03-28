import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import os
import sys
import re


def parse_name(filename: str):
    pattern = re.compile(r"(\d+)MiB")
    match = re.search(pattern, filename)
    if match:
        return int(match.group(1))


def read_data(filepath: str):
    data = pd.read_csv(filepath, header=None, names=["filesize", "rate"])[1:]
    data["filesize"] = data["filesize"].str.extract(r"\d+_(\d+)MB\.data").astype(int)
    data["rate"] = data["rate"].astype(float)
    return data


# parallelism = sys.argv[1]
# ディレクトリとファイルの設定
directory = "../result/buf_capability/journal"
files = os.listdir(directory)
files = sorted(files, key=lambda x: parse_name(x))

# 描画データの準備
plt.figure(figsize=(10, 6))

for file in files:
    # if file.startswith("match") or file.startswith("reten"):
    if file.startswith("match"):
        # print(file)
        filepath = os.path.join(directory, file)
        data = read_data(filepath)

        plt.plot(
            data["filesize"], data["rate"], marker="o", label=f"{parse_name(file)} MiB"
        )

# グラフの設定
plt.xlabel("Size of Original File [MB]", fontsize=18)
plt.xscale("log", base=2)
plt.xticks([2**x for x in range(8)], [str(2**x) for x in range(8)])
plt.ylabel("Match Rate", fontsize=18)
plt.yticks(np.arange(0, 1.1, 0.2))
# plt.title(
#     f"Valiation of Match Rate for Different Ring Buffer Sizes (DoP = {parallelism})",
#     fontsize=16,
# )

plt.legend(title="Ring Buffer Size", fontsize=12)
plt.grid(True)
plt.tight_layout()

# グラフを表示
plt.savefig("./img/buf_capability_seek.png")
