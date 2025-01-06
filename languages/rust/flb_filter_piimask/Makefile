.PHONY: c
c:
	gcc -shared -o bin/pii.so -fPIC -Wall pii.c

.PHONY: rust
rust:
	cargo build --target wasm32-unknown-unknown --release

