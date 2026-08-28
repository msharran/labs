#!/usr/bin/env bash
# Idempotent Cloud Agent bootstrap for this polyglot labs repository.
# The default base image already provides Go, Rust, Python, Node, Java and
# the C toolchain, so this only adds what is missing and warms the demo apps.
set -euo pipefail

# uv: Python package/venv manager used by the python/ experiments
# (e.g. python/realtime-chat pins Python >=3.14 via uv).
if ! command -v uv >/dev/null 2>&1; then
  curl -LsSf https://astral.sh/uv/install.sh | sh
fi
export PATH="$HOME/.local/bin:$PATH"
uv --version

# Warm the flagship Go web app (Fiber, served on :3000).
if [ -d go/go-chat ]; then
  (cd go/go-chat && go mod download)
fi

# Prepare the realtime-chat WebSocket server (served on :8765).
if [ -d python/realtime-chat ]; then
  (cd python/realtime-chat && uv sync)
fi

echo "Cloud Agent install complete."
