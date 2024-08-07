#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void print_usage(const char *progname);
int parse_flag(const char *flag);

int main(int argc, char **argv) {
	if (argc != 3) {
		print_usage(argv[0]);
		return EXIT_FAILURE;
	}

	const char *filepath = argv[1];
	const int is_encrypt = parse_flag(argv[2]);
	if (is_encrypt == -1) {
		print_usage(argv[0]);
		return EXIT_FAILURE;
	}
	printf("Filepath: %s\n", filepath);
}

void print_usage(const char *progname) {
	printf("Usage: %s <filepath> <-d | -e>\n", progname);
}

int parse_flag(const char *flag) {
	if (strcmp(flag, "-e") == 0) {
		return 1;
	} else if (strcmp(flag, "-d") == 0) {
		return 0;
	} else {
		return -1;
	}
}
