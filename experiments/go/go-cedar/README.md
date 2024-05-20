## Install cedar cli

Debian deps

```bash
sudo apt-get update
sudo apt install build-essential -y
```

Installation for all unix based OS

```bash
apt update && apt install build-essential
curl https://sh.rustup.rs -sSf | sh
source "$HOME/.cargo/env"
cargo install cedar-policy-cli --locked
cedar
```
