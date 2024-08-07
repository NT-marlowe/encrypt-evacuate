CC = clang
CFLAGS = -Wall -Wextra -Wconversion -Wdouble-promotion \
     -Wno-unused-parameter -Wno-unused-function -Wno-sign-conversion \
     -fsanitize=undefined -fsanitize-trap 

TARGET = my_simple_ransomware

all: $(TARGET)

$(TARGET): my_simple_ransomware.c
	$(CC) $(CFLAGS) -o $@ $<

clean:
	rm -f $(TARGET)
