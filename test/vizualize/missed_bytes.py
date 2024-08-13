import sys
import matplotlib.pyplot as plt
import numpy as np

BYTES_PER_VALUE = 4


def find_missing_bytes(recovered_file_path: str, byte_size: int) -> list[bool]:
    N = byte_size // BYTES_PER_VALUE
    result = [False] * N
    with open(recovered_file_path, "rb") as f:
        for i in range(1, N + 1):
            # f.seek((i - 1) * BYTES_PER_VALUE)
            val = int.from_bytes(f.read(BYTES_PER_VALUE), byteorder="little")
            result[val - 1] = True

        return result


def viz_boolean_array(arr: list[bool]):
    length = len(arr)
    yoko_length = 1000
    tate_length = length // yoko_length
    # yoko_length = length // tate_length
    arr_2d = np.array(arr).reshape(tate_length, yoko_length)
    # plt.figure(figsize=(20, 20))
    plt.imshow(arr_2d, cmap="gray_r", aspect="auto")

    plt.xticks(
        ticks=np.arange(0, yoko_length, 100),
        labels=np.arange(1, yoko_length + 1, 100),
    )
    plt.yticks(
        ticks=np.arange(0, tate_length, 250), labels=np.arange(0, tate_length, 250)
    )

    # plt.title("Captured Bytes")
    plt.savefig("captured_bytes.png")


recovered_file_path = sys.argv[1]
byte_size = int(sys.argv[2])

viz_boolean_array(find_missing_bytes(recovered_file_path, byte_size))
