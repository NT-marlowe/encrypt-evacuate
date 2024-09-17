data_dir = "../result/ebpf_time"


def parse_line(line: str) -> int:
    time_str = line.split(" ")[-1]
    try:
        time_int = int(time_str)
        return time_int
    except ValueError:
        return 0


def read_data(filepath: str) -> list[int]:
    with open(filepath, "r") as f:
        lines = [l for l in f.readlines() if "time" in l]
        ret = []
        for line in lines:
            line = line.strip("\n")
            time_str = line.split(" ")[-1]
            ret.append(parse_line(line))

        return ret


# caclulate the average and std of the list for each file in data_dir
def calc_avg_std(data_dir: str) -> dict[str, tuple[int, float, float]]:
    import os
    import numpy as np

    ret = {}
    for filename in sorted(os.listdir(data_dir)):
        if "stats" in filename:
            continue
        filepath = f"{data_dir}/{filename}"
        data = read_data(filepath)
        avg = np.mean(data)
        std = np.std(data)
        ret[filename] = (len(data), avg, std)

    return ret


results = calc_avg_std(data_dir)
for key in results:
    elem = results[key]
    print(f"{key}, \t{elem[0]}, \t{elem[1]:.2f}, \t{elem[2]:.2f}")
