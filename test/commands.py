import subprocess
import sys

K = 1000
M = K**2
G = K**3


def generate_random_file(filepath: str, byte_size: int):
    try:
        subprocess.run(["openssl", "rand", "-out", filepath, str(byte_size)])
        print(f"Random file '{filepath}' generated successfully.")
    except FileNotFoundError:
        print("openssl command not found. Please make sure it is installed.")


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
        filesize *= 10
        idx += 1


if __name__ == "__main__":
    generate_files()
