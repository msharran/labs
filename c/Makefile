CC=clang

TARGET=$(shell find . -name "*.c")

# Usage: make <filename>.c
# Example: make hello.c
# Output: bin/hello
.PHONY: $(TARGET)
$(TARGET):
	$(CC) -o bin/$(basename $@) $@
