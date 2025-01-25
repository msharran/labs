#!/bin/bash

run_cmd() {
    eval $@
    if $? -ne 0; then 
        exit 1
    fi
}

echo "*Installing dependencies"
sudo dnf -y update
sudo dnf -y groupinstall "Development Tools"
sudo dnf install -y git\
    bat\
    fzf\
    zoxide\
    eza\
    ripgrep\
    neovim\
    make\
    curl\
    zig\
    golang\
    stow\
    fish\
    btop\
    redis

echo "*Installing dotfiles"
git clone https://github.com/msharran/.dotfiles /home/msharran/.dotfiles
cd /home/msharran/.dotfiles
make

echo "*Installing starship"
curl -sS https://starship.rs/install.sh | sh

echo "*Setting fish as default shell"
sudo chsh -s /usr/bin/fish msharran
