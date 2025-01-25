#!/bin/bash

set -e

echo "*Installing dependencies"
sudo dnf -y update
# sudo dnf -y group install workstation-product-environment
sudo dnf install -y \
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
    btop\
    redis

ZIG_VERSION=0.14.0-dev.2606+b039a8b61
echo "*Installing zig $ZIG_VERSION"

pushd $HOME/.local
    wget https://ziglang.org/builds/zig-linux-aarch64-${ZIG_VERSION}.tar.xz 
    tar -xf zig-linux-aarch64-${ZIG_VERSION}.tar.xz 
    sudo ln -sf $HOME/.local/zig-linux-aarch64-${ZIG_VERSION}/zig /usr/local/bin/zig
    rm zig-linux-aarch64-${ZIG_VERSION}.tar.xz 
    zig version
popd

echo "*Installing dotfiles"
if [[ -d /home/msharran/.dotfiles ]]; then
    echo "* Removing existing dotfiles"
    rm -rf /home/msharran/.dotfiles
fi
git clone https://github.com/msharran/.dotfiles /home/msharran/.dotfiles
cd /home/msharran/.dotfiles
make

echo "*Installing starship"
curl -sS https://starship.rs/install.sh | sh

echo "*Setting fish as default shell"
sudo chsh -s /usr/bin/fish msharran
