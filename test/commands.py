import subprocess
import sys
import os

K = 1000
M = K**2
G = K**3


def generate_random_file(filepath: str, byte_size: int):
    try:
        subprocess.run(["openssl", "rand", "-out", filepath, str(byte_size)])
        print(f"Random file '{filepath}' generated successfully.")
    except FileNotFoundError:
        print("openssl command not found. Please make sure it is installed.")


def write_lorem(filepath: str, byte_size: int):
    with open("lorem_all.txt", "r") as f_lorem:
        string = f_lorem.read()
        with open(filepath, "wb") as f:
            for _ in range(byte_size // len(string)):
                f.write(string.encode("utf-8"))
            f.write(string[: byte_size % len(string)].encode("utf-8"))


def get_suffix(byte_size: int) -> tuple[int, str]:
    if byte_size < K:
        return 1, "B"
    if byte_size < M:
        return K, "KB"
    if byte_size < G:
        return M, "MB"
    return G, "GB"


def get_filename(byte_size: int):
    base, suffix = get_suffix(byte_size)
    return f"{byte_size // base}{suffix}.data"


# 使用例
def generate_files(init_size: int = K):
    filesize = init_size
    idx = 0
    while filesize < G:
        generate_random_file(f"./data/{idx:02d}_{get_filename(filesize)}", filesize)
        # write_lorem(f"./data/{idx:02d}_{get_filename(filesize)}", filesize)
        filesize *= 10
        idx += 1


def calculate_recovery_rate(original_file_path, recovered_file_path):
    with open(original_file_path, "rb") as f:
        original_content = f.read()

    with open(recovered_file_path, "rb") as f:
        recovered_content = f.read()

    # Ensure both files are of the same length for comparison
    length = min(len(original_content), len(recovered_content))
    match_count = sum(
        1 for i in range(length) if original_content[i] == recovered_content[i]
    )

    # Calculate recovery rate
    recovery_rate = (match_count / len(original_content)) * 100
    return recovery_rate


from rapidfuzz import fuzz
import Levenshtein


def calculate_dist_and_ratio(original_file_path, recovered_file_path):
    with open(original_file_path, "r") as f:
        original_content = f.read()

    with open(recovered_file_path, "r") as f:
        recovered_content = f.read()

    # Calculate Levenshtein distance
    distance = Levenshtein.distance(original_content, recovered_content)

    # Calculate the similarity
    similarity = 1 - (distance / max(len(original_content), len(recovered_content)))

    # Calculate recovery rate as a percentage
    recovery_rate = similarity * 100
    partial_ratio = fuzz.partial_ratio(original_content, recovered_content)
    return recovery_rate, partial_ratio


def calculate_retention_rate(original_file_path, recovered_file_path):
    original_size = os.path.getsize(original_file_path)
    recovered_size = os.path.getsize(recovered_file_path)
    return recovered_size / original_size


if __name__ == "__main__":
    subcommand = sys.argv[1]
    if subcommand == "gen":
        generate_files()
        sys.exit(0)

    original_file_path = sys.argv[2]
    recovered_file_path = sys.argv[3]

    if subcommand == "dist":
        # recovery_rate = calculate_recovery_rate(original_file_path, recovered_file_path)
        recovery_rate, partial_ratio = calculate_dist_and_ratio(
            original_file_path, recovered_file_path
        )
        print(f"{recovery_rate:.2f}, {partial_ratio}")

    elif subcommand == "match":
        recovery_rate = calculate_recovery_rate(original_file_path, recovered_file_path)
        print(f"{recovery_rate:.2f}")

    elif subcommand == "reten":
        retention_rate = calculate_retention_rate(
            original_file_path, recovered_file_path
        )
        print(f"{retention_rate:.3f}")

# generate_files()
