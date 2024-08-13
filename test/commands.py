import sys
import os
from rapidfuzz import fuzz
from subcommands.gen import generate_files
import Levenshtein


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
    if subcommand == "gen-rand":
        generate_files()
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
        print(f"{recovery_rate:.2f}, {partial_ratio}")

    elif subcommand == "match":
        recovery_rate = calculate_recovery_rate(original_file_path, recovered_file_path)
        print(f"{recovery_rate:.2f}")

    elif subcommand == "reten":
        retention_rate = calculate_retention_rate(
            original_file_path, recovered_file_path
        )
        print(f"{retention_rate:.3f}")

    else:
        print("Invalid subcommand.")
        sys.exit(1)
