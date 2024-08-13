import matplotlib.pyplot as plt
import re
import sys

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


data_path = sys.argv[1]
filesize_list, ratio_list = read_csv(f"{data_path}")
print(filesize_list, ratio_list)
