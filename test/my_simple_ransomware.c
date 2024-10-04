#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "./cipher_funcs.h"

void print_usage(const char *progname);

int main(int argc, char **argv) {
	if (argc != 2) {
		print_usage(argv[0]);
		return EXIT_FAILURE;
	}

	const char *filepath = argv[1];

	unsigned char key[32] = {0};
	unsigned char iv[16]  = {0};
	// Generate a random key and IV
	// RAND_bytes(key, sizeof(key));
	// RAND_bytes(iv, sizeof(iv));

	encrypt_file(filepath, key, iv);
}

void print_usage(const char *progname) {
	printf("Usage: %s <filepath>\n", progname);
}
