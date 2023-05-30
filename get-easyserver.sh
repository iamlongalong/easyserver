#!/bin/bash

set -e

# Determine the OS type and architecture
OS_TYPE=$(uname -s)
OS_ARCH=$(uname -m)

LATEST_URL="https://api.github.com/repos/iamlongalong/easyserver/releases/latest"

# Determine the download URL based on the OS type and architecture
case $OS_TYPE in
Linux)
  if [ "$OS_ARCH" = "x86_64" ]; then
    DOWNLOAD_URL=$(curl -s "$LATEST_URL" | grep "browser_download_url.*linux-x86_64" | cut -d '"' -f 4)
    # DOWNLOAD_URL="https://github.com/iamlongalong/easyserver/releases/latest/download/easyserver-linux-x86_64"
  elif [ "$OS_ARCH" = "aarch64" ]; then
    DOWNLOAD_URL=$(curl -s "$LATEST_URL" | grep "browser_download_url.*linux-arm64" | cut -d '"' -f 4)
    # DOWNLOAD_URL="https://github.com/iamlongalong/easyserver/releases/latest/download/easyserver-linux-arm64"
  else
    echo "😭 Unsupported CPU architecture: $OS_ARCH"
    exit 1
  fi
  ;;
Darwin)
  if [ "$(uname -m)" = "x86_64" ]; then
    # x86-64 Mac
    DOWNLOAD_URL=$(curl -s "$LATEST_URL" | grep "browser_download_url.*darwin-x86_64" | cut -d '"' -f 4)
    # DOWNLOAD_URL="https://github.com/iamlongalong/easyserver/releases/latest/download/easyserver-darwin-x86_64"
  elif [ "$(uname -m)" = "arm64" ]; then
    # M1 Mac
    DOWNLOAD_URL=$(curl -s "$LATEST_URL" | grep "browser_download_url.*darwin-arm64" | cut -d '"' -f 4)
    # DOWNLOAD_URL="https://github.com/iamlongalong/easyserver/releases/latest/download/easyserver-darwin-arm64"
  else
    echo "😭 Unsupported Mac architecture: $(uname -m)"
    exit 1
  fi
  ;;
CYGWIN* | MINGW*)
  if [ "$OS_ARCH" = "x86_64" ]; then
    DOWNLOAD_URL=$(curl -s "$LATEST_URL" | grep "browser_download_url.*windows-x86_64" | cut -d '"' -f 4)
    # DOWNLOAD_URL="https://github.com/iamlongalong/easyserver/releases/latest/download/easyserver-windows-x86_64"
  else
    echo "😭 Unsupported CPU architecture: $OS_ARCH"
    exit 1
  fi
  ;;
*)
  echo "😭 Unsupported operating system: $OS_TYPE"
  exit 1
  ;;
esac

if [ -z "$DOWNLOAD_URL" ]; then
  echo "😭 无法找到适合当前系统的 easyserver 二进制文件"
  exit 1
fi

# Download and extract the latest release
TMP_DIR=$(mktemp -d)

# if [ -n "$DOWNLOAD_URL" ]; then
curl -L -o "$TMP_DIR"/easyserver "$DOWNLOAD_URL"
chmod +x "$TMP_DIR"/easyserver
# fi

# Move the easyserver binary to /usr/local/bin
sudo mv "$TMP_DIR"/easyserver /usr/local/bin/easyserver

# Clean up
rm -rf "$TMP_DIR"

echo
echo "🎉🎉🎉 easyserver 已更新到最新版本！ 🎉🎉🎉"
echo

/usr/local/bin/easyserver -h
