import pandas as pd
import matplotlib.pyplot as plt
import re
import sys

operations = [
    "rd.Read",
    "binary.Read",
    "MakeItem",
    "minHeapSort",
    "file.Write",
]
parallelism = 4

# ファイルからデータを読み込む
data = {operation: [] for operation in operations}
with open(sys.argv[1], "r") as file:
    for line in file:
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

print()
for ops in operations:
    print(f"{ops}: {len(data[ops])} samples")

# データフレームに変換
df = pd.DataFrame(dict([(k, pd.Series(v)) for k, v in data.items()]))

# 箱ひげ図の描画
"""
DataFrameのカラムから箱ひげ図を作成します。 
箱ひげ図は、数値データのグループを四分位値でグラフ化する方法です。 
箱はデータのQ1からQ3の四分位値から広がり、中央値 (Q2)に線が引かれる。 
ひげは、データの範囲を示すためにボックスの端から伸びます。 
デフォルトでは、ボックスの端から1.5 * IQR (IQR = Q3 - Q1)を超えない範囲で広がり、
その区間内で最も遠いデータポイントで終わります。 外れ値は，個別の点としてプロットされる．
"""
plt.figure(figsize=(12, 7))
df.boxplot()
plt.ylabel("Time (us)")
plt.title(f"Processing Time, parallelism = {parallelism} (> 1000us are ignored)")
plt.xticks(rotation=30)
plt.grid(True)
plt.savefig(f"./img/hakohige_parallel_{parallelism}.png")
