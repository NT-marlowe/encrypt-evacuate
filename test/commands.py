import sys
from subcommands.gen import generate_random_files, generate_sequential_file
from subcommands.calc import (
    calculate_match_rate,
    calculate_retention_rate,
)


if __name__ == "__main__":
    subcommand = sys.argv[1]
    if subcommand == "gen-rand":
        generate_random_files()
        sys.exit(0)

    if subcommand == "gen-seq":
        byte_size = int(sys.argv[2])
        generate_sequential_file(byte_size)
        sys.exit(0)

    original_file_path = sys.argv[2]
    sheltered_file_path = sys.argv[3]

    if subcommand == "match":
        recovery_rate = calculate_match_rate(original_file_path, sheltered_file_path)
        print(f"{recovery_rate:.3f}")

    elif subcommand == "reten":
        retention_rate = calculate_retention_rate(
            original_file_path, sheltered_file_path
        )
        print(f"{retention_rate:.3f}")

    else:
        print("Invalid subcommand.")
        sys.exit(1)
