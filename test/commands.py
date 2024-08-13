import sys
import os
from rapidfuzz import fuzz
from subcommands.gen import generate_random_files
from subcommands.calc import (
    calculate_dist_and_ratio,
    calculate_recovery_rate,
    calculate_retention_rate,
)


if __name__ == "__main__":
    subcommand = sys.argv[1]
    if subcommand == "gen-rand":
        generate_random_files()
        sys.exit(0)

    if subcommand == "gen-seq":
        byte_size = int(sys.argv[2])

    original_file_path = sys.argv[2]
    recovered_file_path = sys.argv[3]

    if subcommand == "dist":
        # recovery_rate = calculate_recovery_rate(original_file_path, recovered_file_path)
        recovery_rate, partial_ratio = calculate_dist_and_ratio(
            original_file_path, recovered_file_path
        )
        print(f"{recovery_rate:.3f}, {partial_ratio}")

    elif subcommand == "match":
        recovery_rate = calculate_recovery_rate(original_file_path, recovered_file_path)
        print(f"{recovery_rate:.3f}")

    elif subcommand == "reten":
        retention_rate = calculate_retention_rate(
            original_file_path, recovered_file_path
        )
        print(f"{retention_rate:.3f}")

    else:
        print("Invalid subcommand.")
        sys.exit(1)
