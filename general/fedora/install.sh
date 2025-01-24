#!/bin/bash
sudo dnf -y groupinstall "Development Tools"
sudo dnf install -y git bat fzf zoxide eza ripgrep neovim make curl zig golang stow fish

# Install dotfiles
git clone https://github.com/msharran/.dotfiles /home/msharran/.dotfiles
cd /home/msharran/.dotfiles
make

# Configure shell 
curl -sS https://starship.rs/install.sh | sh
sudo chsh -s /usr/bin/fish msharran
