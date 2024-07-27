#include <openssl/evp.h>
#include <stdio.h>
#include <string.h>

int main(void) {
	unsigned char *plaintext = (unsigned char *)"This is a secret message.";
	unsigned char ciphertext[128];
	unsigned char decryptedtext[128];

	EVP_CIPHER_CTX *ctx;
	int len;
	int ciphertext_len;
	int plaintext_len = strlen((char *)plaintext);

	unsigned char key[32];
	unsigned char iv[16];

	// コンテキストの作成
	if (!(ctx = EVP_CIPHER_CTX_new()))
		return -1;

	// 暗号化の初期化
	if (1 != EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv))
		return -1;

	// データの暗号化
	if (1 != EVP_EncryptUpdate(ctx, ciphertext, &len, plaintext, plaintext_len))
		return -1;

	printf("ciphertext: %s\n", ciphertext);

	return 0;
}
