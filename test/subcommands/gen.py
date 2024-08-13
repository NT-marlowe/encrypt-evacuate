import subprocess

K = 1000
M = K**2
G = K**3


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


def generate_random_files(init_size: int = K):
    filesize = init_size
    idx = 0
    while filesize < G:
        generate_random_file(f"./data/{idx:02d}_{get_filename(filesize)}", filesize)
        # write_lorem(f"./data/{idx:02d}_{get_filename(filesize)}", filesize)
        filesize *= 10
        idx += 1
