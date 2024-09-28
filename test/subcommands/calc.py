from rapidfuzz import fuzz
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


# def calculate_dist_and_ratio(original_file_path, sheltered_file_path):
#     with open(original_file_path, "r") as f:
#         original_content = f.read()

#     with open(sheltered_file_path, "r") as f:
#         recovered_content = f.read()

#     # Calculate Levenshtein distance
#     distance = Levenshtein.distance(original_content, recovered_content)

#     # Calculate the similarity
#     similarity = 1 - (distance / max(len(original_content), len(recovered_content)))

#     # Calculate recovery rate as a percentage
#     recovery_rate = similarity
#     partial_ratio = fuzz.partial_ratio(original_content, recovered_content)
#     return recovery_rate, partial_ratio


def calculate_retention_rate(original_file_path, sheltered_file_path):
    original_size = os.path.getsize(original_file_path)
    recovered_size = os.path.getsize(sheltered_file_path)
    return recovered_size / original_size
