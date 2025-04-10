import matplotlib.pyplot as plt
import re
import numpy as np
import sys


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


def plot_graph(metrics: str, file1, file2, parallelsim: int, size_type: str):
    # def plot_graph(metrics: str, filename: str, parallelsim: int, size_type: str):
    # _, values = read_csv(filename)
    _, list1 = read_csv(file1)
    _, list2 = read_csv(file2)

    print(list1)
    print(list2)

    if size_type == "exp":
        labels = ["1K", "10K", "100K", "1M", "10M", "100M"]
    elif size_type == "inc":
        labels = ["1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M", "10M"]
    else:
        raise ValueError("Invalid size type")

    x = np.arange(len(labels))
    width = 0.2

    fig, ax = plt.subplots()
    # ax.bar(x - width / 2, list1, width, label="Preliminary Exp.")
    # ax.bar(x + width / 2, list2, width, label="Main Exp.")
    ax.bar(x - width / 2, list1, width, label="Prior Work [15]")
    ax.bar(x + width / 2, list2, width, label="Proposed")
    # ax.bar(x, values, width)

    # if metrics == "reten":
    #     ax.set_title(f"Retention Rates by File Size, p = {parallesim:02}")
    # elif metrics == "match":
    #     ax.set_title(f"Match Rates by File Size, p = {parallesim:02}")

    ax.set_xlabel("Size of Original File [Byte]")
    ax.set_ylabel("Match Rate")
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    # plot these two lists
    # plt.savefig(f"./img/seqential_vs_parallel_{metrics}_p{parallelsim}_{size_type}.png")
    plt.savefig(f"./img/{metrics}_p{parallelsim}_{size_type}.png")


# plot_graph_exp()

# plot_graph_incremental()

if __name__ == "__main__":
    parallesim = int(sys.argv[1])

    file1, file2 = sys.argv[2], sys.argv[3]
    # filename = sys.argv[2]

    size_type = sys.argv[4]

    # plot_graph_inc("reten")
    plot_graph("match", file1, file2, parallesim, size_type)
    # plot_graph("match", filename, parallesim, size_type)
