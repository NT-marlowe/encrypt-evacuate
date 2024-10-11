import matplotlib.pyplot as plt
import re
import numpy as np

data_dir_exp = "../result/exponential"
data_dir_inc = "../result/incremental"


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


def plot_graph_exp():
    filesize_list, retention_list = read_csv(f"{data_dir_exp}/retention.csv")
    filesize_list, match_list = read_csv(f"{data_dir_exp}/match.csv")

    labels = ["1K", "10K", "100K", "1M", "10M", "100M"]
    x = np.arange(len(labels))
    width = 0.2

    fig, ax = plt.subplots()
    ax.bar(x - width / 2, retention_list, width, label="Retention Rate")
    ax.bar(x + width / 2, match_list, width, label="Match Rate")

    ax.set_title("Retention and Match Rates by File Size")

    ax.set_xlabel("Size of Original File [Byte]")
    ax.set_ylabel("Rate")
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    # plot these two lists
    plt.savefig(f"./img/retention_match_exp.png")


def plot_graph(metrics: str):
    _, match_list = read_csv(f"{data_dir_inc}/{metrics}.csv")
    _, match_list_p4 = read_csv(f"{data_dir_inc}/{metrics}_p4.csv")

    labels = ["1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M", "10M"]
    x = np.arange(len(labels))
    width = 0.2

    fig, ax = plt.subplots()
    ax.bar(x - width / 2, match_list, width, label="Seqential")
    ax.bar(x + width / 2, match_list_p4, width, label="Parallel (p = 4)")

    if metrics == "retention":
        ax.set_title("Retention Rates by File Size")
    elif metrics == "match":
        ax.set_title("Match Rates by File Size")

    ax.set_xlabel("Size of Original File [Byte]")
    ax.set_ylabel("Rate")
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    # plot these two lists
    plt.savefig(f"./img/{metrics}_seq_paralell.png")


def plot_graph_incremental():
    _, match_list = read_csv(f"{data_dir_inc}/match.csv")
    _, match_list_p4 = read_csv(f"{data_dir_inc}/match_p4.csv")

    # _, retention_list = read_csv(f"{data_dir_inc}/retention.csv")
    # _, retention_p4_list = read_csv(f"{data_dir_inc}/retention_p4.csv")

    labels = ["1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M", "10M"]
    x = np.arange(len(labels))
    width = 0.2

    fig, ax = plt.subplots()
    # ax.bar(x - width / 2, retention_list, width, label="Retention Rate (seq.)")
    # ax.bar(x + width / 2, retention_p4_list, width, label="Retention Rate (p = 4)")
    ax.bar(x - width / 2, match_list, width, label="Match Rate (seq.)")
    ax.bar(x + width / 2, match_list_p4, width, label="Match Rate (p = 4)")

    # ax.set_title("Retention Rates by File Size")
    ax.set_title("Match Rates by File Size")

    ax.set_xlabel("Size of Original File [Byte]")
    ax.set_ylabel("Rate")
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    # plot these two lists
    # plt.savefig(f"./img/retention_match_inc.png")
    # plt.savefig(f"./img/retention.png")
    plt.savefig(f"./img/match.png")


# plot_graph_exp()

# plot_graph_incremental()
plot_graph("retention")
plot_graph("match")
