CC=clang

TARGET=$(shell find . -name "*.c")

.PHONY: $(TARGET)
$(TARGET):
	$(CC) -o bin/$(basename $@) $@
