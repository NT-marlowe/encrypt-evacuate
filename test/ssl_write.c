#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <openssl/ssl.h>
#include <openssl/err.h>

void initialize_ssl() {
	SSL_load_error_strings();
	OpenSSL_add_ssl_algorithms();
}

void cleanup_ssl() {
	EVP_cleanup();
}

SSL_CTX *create_context() {
	const SSL_METHOD *method;
	SSL_CTX *ctx;

	method = SSLv23_client_method();

	ctx = SSL_CTX_new(method);
	if (!ctx) {
		perror("Unable to create SSL context");
		ERR_print_errors_fp(stderr);
		exit(EXIT_FAILURE);
	}

	return ctx;
}

void configure_context(SSL_CTX *ctx) {
	SSL_CTX_set_ecdh_auto(ctx, 1);
}

int main(int argc, char **argv) {
	const char *hostname = "www.example.com";
	const char *portnum  = "443";
	int sock;
	struct sockaddr_in server_addr;
	SSL_CTX *ctx;
	SSL *ssl;

	initialize_ssl();

	ctx = create_context();
	configure_context(ctx);

	sock = socket(AF_INET, SOCK_STREAM, 0);
	if (sock < 0) {
		perror("Unable to create socket");
		exit(EXIT_FAILURE);
	}

	memset(&server_addr, 0, sizeof(server_addr));
	server_addr.sin_family = AF_INET;
	server_addr.sin_port   = htons(atoi(portnum));

	if (inet_pton(AF_INET, "93.184.216.34", &server_addr.sin_addr) <=
		0) { // Example.com IP
		perror("Invalid address/ Address not supported");
		exit(EXIT_FAILURE);
	}

	if (connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) <
		0) {
		perror("Connection Failed");
		close(sock);
		exit(EXIT_FAILURE);
	}

	ssl = SSL_new(ctx);
	SSL_set_fd(ssl, sock);

	if (SSL_connect(ssl) <= 0) {
		ERR_print_errors_fp(stderr);
	} else {
		const char *msg = "Hello, SSL!";
		SSL_write(ssl, msg, strlen(msg));
		printf("Sent: %s\n", msg);
	}

	SSL_free(ssl);
	close(sock);
	SSL_CTX_free(ctx);
	cleanup_ssl();

	return 0;
}
