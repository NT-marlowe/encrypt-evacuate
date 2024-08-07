#pragma once

#define BUFFER_SIZE 4096

void handle_errors(void);
void encrypt_file(const char *input_filepath, const unsigned char *key,
	const unsigned char *iv);
