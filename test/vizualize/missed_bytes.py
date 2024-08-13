import sys

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


recovered_file_path = sys.argv[1]
byte_size = int(sys.argv[2])
result = find_missing_bytes(recovered_file_path, byte_size)
print(f"count of True: {sum(result)}")
