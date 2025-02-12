#!/bin/bash

set -e

echo "*Installing dependencies"
sudo apt-get -y update
sudo apt-get install -y \
    wget\
    git\
    bat\
    fzf\
    zoxide\
    eza\
    ripgrep\
    neovim\
    make\
    curl\
    golang\
    stow\
    fish\
    btop

ZIG_VERSION=0.14.0-dev.2606+b039a8b61
echo "*Installing zig $ZIG_VERSION"

pushd $HOME/.local
    wget https://ziglang.org/builds/zig-linux-aarch64-${ZIG_VERSION}.tar.xz
    tar -xf zig-linux-aarch64-${ZIG_VERSION}.tar.xz
    sudo ln -sf $HOME/.local/zig-linux-aarch64-${ZIG_VERSION}/zig /usr/local/bin/zig
    rm zig-linux-aarch64-${ZIG_VERSION}.tar.xz
    zig version
popd

echo "*Installing starship"
curl -sS https://starship.rs/install.sh | sh

echo "*Setting fish as default shell"
sudo chsh -s /usr/bin/fish msharran

echo "*Setting up projects dir"
mkdir -p $HOME/projects/{work,play}

echo "*Installing dotfiles"
if [ ! -d "$HOME/.dotfiles" ]; then
    git clone https://github.com/msharran/.dotfiles $HOME/.dotfiles
    pushd $HOME/.dotfiles
        make
    popd
fi
