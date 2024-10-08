import matplotlib.pyplot as plt
import pandas as pd
import os
import re


def parse_name(filename: str):
    pattern = re.compile(r"(\d+)")
    match = re.search(pattern, filename)
    if match:
        return int(match.group(1))


# ディレクトリとファイルの設定
directory = "../result/buf_capability"
files = os.listdir(directory)
files = sorted(files, key=lambda x: parse_name(x))

# 描画データの準備
plt.figure(figsize=(10, 6))

for file in files:
    # if file.startswith("match") or file.startswith("reten"):
    if file.startswith("match"):
        # print(file)
        filepath = os.path.join(directory, file)
        data = pd.read_csv(filepath, header=None, names=["filesize", "rate"])
        # ファイルサイズをMB単位で表示できるように処理
        data["filesize"] = (
            data["filesize"].str.extract(r"\d+_(\d+)MB\.data").astype(int)
        )

        plt.plot(
            data["filesize"], data["rate"], marker="o", label=f"{parse_name(file)} MiB"
        )

# グラフの設定
plt.xlabel("Size of Original File [MB]", fontsize=16)
plt.xscale("log", base=2)
plt.xticks([2**x for x in range(8)], [str(2**x) for x in range(8)])
plt.ylabel("Rate", fontsize=16)
plt.title(
    "Valiation of Match Rate with File Size for Different Ring Buffer Sizes",
    fontsize=16,
)
plt.legend(title="Ring Buffer Size", fontsize=10)
plt.grid(True)
plt.tight_layout()

# グラフを表示
plt.savefig("./img/buf_capability.png")
