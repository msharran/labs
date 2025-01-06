FROM mcr.microsoft.com/devcontainers/base:bookworm

# Note: The DEBIAN_FRONTEND export avoids warnings when you go on to work with your container.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    # && apt-get -y install git

ARG ZIG_VERSION=0.14.0-dev.2606+b039a8b61

RUN cd /opt \
    && wget https://ziglang.org/builds/zig-linux-aarch64-${ZIG_VERSION}.tar.xz \
    && tar -xf zig-linux-aarch64-${ZIG_VERSION}.tar.xz \
    && ln -s /opt/zig-linux-aarch64-${ZIG_VERSION}/zig /usr/local/bin/zig \
    && rm zig-linux-aarch64-${ZIG_VERSION}.tar.xz \
    && zig version

# TODO Install nvim from source, no arm64 linux binary available
# Ref: https://github.com/neovim/neovim/blob/master/BUILD.md
RUN cd /opt \
    && apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && sudo apt-get install ninja-build gettext cmake unzip curl build-essential

WORKDIR /workspace
