#include "cipher_funcs.h"

#include <errno.h>
#include <openssl/aes.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void handle_errors(void) {
	ERR_print_errors_fp(stderr);
	abort();
}

void encrypt_file(const char *input_filepath, const unsigned char *key,
	const unsigned char *iv) {
	FILE *input_file = fopen(input_filepath, "rb");
	if (!input_file) {
		perror("fopen");
		exit(EXIT_FAILURE);
	}

	// get a filepath `input_filepath + ".enc"`
	char output_filepath[100];
	strcpy(output_filepath, input_filepath);
	strcat(output_filepath, ".enc");
	// const char *output_filepath = strcat(input_filepath, ".enc");

	FILE *output_file = fopen(output_filepath, "wb");
	if (!output_file) {
		perror("fopen");
		fclose(input_file);
		exit(EXIT_FAILURE);
	}

	EVP_CIPHER_CTX *ctx = EVP_CIPHER_CTX_new();
	if (!ctx) {
		handle_errors();
	}

	if (EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv) != 1) {
		handle_errors();
	}

	unsigned char buffer[BUFFER_SIZE];
	size_t bytes_read = 0;
	int bytes_written = 0;
	// int accumulated_bytes_read = 0;
	while ((bytes_read = fread(buffer, 1, BUFFER_SIZE, input_file)) > 0) {
		// accumulated_bytes_read += bytes_read;
		printf("read bytes : %d\n", bytes_read);
		if (EVP_EncryptUpdate(
				ctx, buffer, &bytes_written, buffer, (int)bytes_read) != 1) {
			handle_errors();
		}
		fwrite(buffer, 1, bytes_written, output_file);
	}

	if (EVP_EncryptFinal_ex(ctx, buffer, &bytes_written) != 1) {
		handle_errors();
	}
	fwrite(buffer, 1, bytes_written, output_file);

	EVP_CIPHER_CTX_free(ctx);
	fclose(input_file);
	fclose(output_file);
}
