import os


def calculate_match_rate(original_file_path, sheltered_file_path):
    with open(original_file_path, "rb") as f:
        original_content = f.read()

    with open(sheltered_file_path, "rb") as f:
        sheltered_content = f.read()

    # Ensure both files are of the same length for comparison
    length = min(len(original_content), len(sheltered_content))
    match_count = sum(
        1 for i in range(length) if original_content[i] == sheltered_content[i]
    )

    # Calculate recovery rate
    match_rate = match_count / len(original_content)
    return match_rate


def calculate_retention_rate(original_file_path, sheltered_file_path):
    original_size = os.path.getsize(original_file_path)
    recovered_size = os.path.getsize(sheltered_file_path)
    return recovered_size / original_size
