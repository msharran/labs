#! /bin/bash
set -e

SHELL=$1

RC_FILE=.bashrc
if [[ $SHELL == "zsh" ]]; then
  RC_FILE="$HOME/.zshrc"
elif [[ $SHELL == "bash" ]]; then
  RC_FILE="$HOME/.bashrc"
else
  echo "Missing argument [SHELL]. List of supported shells [zsh, bash]"
  exit 1
fi

# Create configuration file if it doesn't exist
if [ ! -f "$RC_FILE" ]; then
  echo "$RC_FILE does not exist. creating it"
  touch $RC_FILE
fi

echo "Adding ubuntu funtion to $RC_FILE"

cat <<EOF >> $RC_FILE

### UBUNTU CLI CONFIGURATION

ubuntu() {
  if [[ \$1 == "-h" ]] || [[ \$1 == "--help" ]]; then
    printf "%s\n" \
    "Command Usage:" \
    "  ubuntu start" \
    "  ubuntu stop" \
    "  ubuntu exec bash" \
    "  ubuntu exec ls -l"
  else
    UBUNTU_DIR=\$HOME/play/labs/linux/ubuntu
    if [[ -d \$UBUNTU_DIR ]]; then
      cd \$UBUNTU_DIR 
      ./\$1.sh \${@:2}
    else
      echo "directory not found. clone github.com/msharran/labs into \$HOME/play and try again"
    fi
  fi
}

EOF

echo "Installed ubuntu CLI!"
echo "To use it, open new terminal or run \"source $RC_FILE\"."
printf "%s\n" \
  "Command Usage:" \
  "  ubuntu start" \
  "  ubuntu stop" \
  "  ubuntu exec bash" \
  "  ubuntu exec ls -l"

