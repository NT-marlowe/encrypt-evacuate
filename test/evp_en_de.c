#include <openssl/evp.h>
#include <openssl/aes.h>
#include <openssl/rand.h>

int encrypt(unsigned char *plaintext, int plaintext_len, unsigned char *key,
	unsigned char *iv, unsigned char *ciphertext) {
	EVP_CIPHER_CTX *ctx;
	int len;
	int ciphertext_len;

	// コンテキストの作成
	if (!(ctx = EVP_CIPHER_CTX_new()))
		return -1;

	// 暗号化の初期化
	if (1 != EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv))
		return -1;

	// データの暗号化
	if (1 != EVP_EncryptUpdate(ctx, ciphertext, &len, plaintext, plaintext_len))
		return -1;
	ciphertext_len = len;

	// 最後のブロックを処理
	if (1 != EVP_EncryptFinal_ex(ctx, ciphertext + len, &len))
		return -1;
	ciphertext_len += len;

	// コンテキストのクリーンアップ
	EVP_CIPHER_CTX_free(ctx);

	return ciphertext_len;
}

int decrypt(unsigned char *ciphertext, int ciphertext_len, unsigned char *key,
	unsigned char *iv, unsigned char *plaintext) {
	EVP_CIPHER_CTX *ctx;
	int len;
	int plaintext_len;

	// コンテキストの作成
	if (!(ctx = EVP_CIPHER_CTX_new()))
		return -1;

	// 復号化の初期化
	if (1 != EVP_DecryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv))
		return -1;

	// データの復号化
	if (1 !=
		EVP_DecryptUpdate(ctx, plaintext, &len, ciphertext, ciphertext_len))
		return -1;
	plaintext_len = len;

	// 最後のブロックを処理
	if (1 != EVP_DecryptFinal_ex(ctx, plaintext + len, &len))
		return -1;
	plaintext_len += len;

	// コンテキストのクリーンアップ
	EVP_CIPHER_CTX_free(ctx);

	return plaintext_len;
}

#include <stdio.h>
#include <string.h>

int main(void) {
	// テキスト
	unsigned char *plaintext = (unsigned char *)"This is a secret message.";
	unsigned char ciphertext[128];
	unsigned char decryptedtext[128];

	// 鍵とIV
	unsigned char key[32];
	unsigned char iv[16];

	// 鍵とIVを生成
	if (!RAND_bytes(key, sizeof(key)) || !RAND_bytes(iv, sizeof(iv))) {
		fprintf(stderr, "RAND_bytes error\n");
		return 1;
	}

	// 暗号化
	int ciphertext_len =
		encrypt(plaintext, strlen((char *)plaintext), key, iv, ciphertext);
	if (ciphertext_len < 0) {
		fprintf(stderr, "Encryption error\n");
		return 1;
	}

	// 復号化
	int decryptedtext_len =
		decrypt(ciphertext, ciphertext_len, key, iv, decryptedtext);
	if (decryptedtext_len < 0) {
		fprintf(stderr, "Decryption error\n");
		return 1;
	}

	// 終端文字を追加
	decryptedtext[decryptedtext_len] = '\0';

	// 結果の表示
	printf("Ciphertext is:\n");
	BIO_dump_fp(stdout, (const char *)ciphertext, ciphertext_len);
	printf("Decrypted text is: %s\n", decryptedtext);

	return 0;
}
