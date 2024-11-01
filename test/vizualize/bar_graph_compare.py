import matplotlib.pyplot as plt
import re
import numpy as np
import sys


parallesim = int(sys.argv[1])
data_dir = f"../result/reorder_vs_seek/p{parallesim}"


def suffix_to_unit(suffix: str) -> int:
    if suffix == "K":
        return 1000
    if suffix == "M":
        return 1000**2
    if suffix == "G":
        return 1000**3
    return 1


def get_filesize(filename: str) -> int:
    pattern = re.compile(r"\d+_(\d+)([a-zA-Z]+)B\.data")
    # file pattern: 10_1MB.data
    match = re.search(pattern, filename)
    if match:
        num = int(match.group(1))
        suffix = match.group(2)
        return num * suffix_to_unit(suffix)

    return 0


def read_csv(filepath: str) -> tuple[list[int], list[float]]:
    with open(filepath, "r") as f:
        # skip the first line
        lines = f.readlines()[1:]
        filesize_list = []
        ratio_list = []
        for line in lines:
            filename, ratio = line.strip("\n").split(",")
            filesize_list.append(get_filesize(filename))
            ratio_list.append(float(ratio))

        return filesize_list, ratio_list


def plot_graph_inc(metrics: str):
    _, list_reorder = read_csv(f"{data_dir}/fix_reorder_{metrics}.txt")
    _, list_seek = read_csv(f"{data_dir}/seek_write_{metrics}.txt")

    labels = ["1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M", "10M"]
    x = np.arange(len(labels))
    width = 0.2

    fig, ax = plt.subplots()
    ax.bar(x - width / 2, list_reorder, width, label="Fix Reorder")
    ax.bar(x + width / 2, list_seek, width, label="Seek & Write")

    if metrics == "reten":
        ax.set_title(f"Retention Rates by File Size, p = {parallesim}")
    elif metrics == "match":
        ax.set_title(f"Match Rates by File Size, p = {parallesim}")

    ax.set_xlabel("Size of Original File [Byte]")
    ax.set_ylabel("Rate")
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    # plot these two lists
    plt.savefig(f"./img/reorder_vs_seek_{metrics}_p{parallesim}.png")


# plot_graph_exp()

# plot_graph_incremental()
plot_graph_inc("reten")
plot_graph_inc("match")
