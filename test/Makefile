CC = clang
CFLAGS = -Wall -Wextra -Wconversion -Wdouble-promotion \
     -Wno-unused-parameter -Wno-unused-function -Wno-sign-conversion \
     -fsanitize=undefined -fsanitize-trap 

LDFLAGS = -lcrypto -lssl

TARGET = my_simple_ransomware

all: $(TARGET)

$(TARGET): my_simple_ransomware.c cipher_funcs.c cipher_funcs.h
	$(CC) $(CFLAGS) -o $@ $< cipher_funcs.c $(LDFLAGS)

clean:
	rm -f $(TARGET)
