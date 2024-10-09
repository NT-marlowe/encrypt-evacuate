import pandas as pd
import matplotlib.pyplot as plt
import re
import sys

# ファイルからデータを読み込む
operations = ["rd.Read", "binary.Read", "file.Write", "minHeapSort"]
data = {operation: [] for operation in operations}
parallelism = int(sys.argv[2])

with open(sys.argv[1], "r") as file:
    for line in file:
        # match = re.match(r"(\w+\.\w+):\s([\d\.]+)([a-zµ]*)", line)
        match = re.match(r"([\w\.]+):\s([\d\.]+)([a-zµ]*)", line)

        if match:
            operation, value, unit = match.groups()
            value = float(value)
            # 単位をusに変換
            if unit == "s":
                value *= 1000 * 1000
            elif unit == "ms":
                value *= 1000
            elif unit == "ns":
                value /= 1000

            if operation == "binary.Read":
                value /= parallelism

            if value >= 1000:
                print(f"{operation} took {value} us")
                continue
            if operation in data:
                data[operation].append(value)

# データフレームに変換
df = pd.DataFrame(dict([(k, pd.Series(v)) for k, v in data.items()]))

# 箱ひげ図の描画
plt.figure(figsize=(12, 7))
df.boxplot()
plt.ylabel("Time (us)")
plt.title("Processing Time for Different Operations (> 1000us are ignored)")
plt.xticks(rotation=45)
plt.grid(True)
plt.savefig("./img/hakohige.png")
