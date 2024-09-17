#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

ZIRVA_PORT=""
INSTALL_DIR="/opt/zirva"
LATEST_RELEASE_URL="https://api.github.com/repos/zirvaorg/server/releases/latest"
LATEST_VERSION=$(curl -s $LATEST_RELEASE_URL | grep '"tag_name":' | cut -d '"' -f 4)

echo -e "${BLUE}... zirva server updater ...${NC}"
echo -e "${YELLOW}[info] checking for new version of zirva...${NC}"

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[err] please run this script as root (using sudo).${NC}"
  exit 1
fi

if [ -f "$INSTALL_DIR/zirva" ]; then
  CURRENT_VERSION=$("$INSTALL_DIR/zirva" -v)
else
  echo -e "${RED}[err] zirva is not installed.${NC}"
  exit 1
fi

if [ -z "$LATEST_VERSION" ]; then
  echo -e "${RED}[err] failed to check the latest version.${NC}"
  exit 1
fi

if [ "$CURRENT_VERSION" = "$LATEST_VERSION" ]; then
  echo -e "${GREEN}[success] you are already using the latest version: $CURRENT_VERSION${NC}"
  exit 1
else
  echo -e "${YELLOW}[info] new version available: $LATEST_VERSION. updating...${NC}"
fi

echo -e "${YELLOW}[info] checking if zirva is running...${NC}"
ZIRVA_PID=$(pgrep -f "zirva")
if [ -z "$ZIRVA_PID" ]; then
  echo -e "${GREEN}[success] zirva is not running.${NC}"
else
  ZIRVA_PORT=$(ps -fp $ZIRVA_PID | awk '/zirva/ {for(i=1;i<=NF;i++) if($i=="-p") print $(i+1)}' | sort -u | tr -d '\n')
  echo -e "${YELLOW}[info] zirva is running on port $ZIRVA_PORT.${NC}"
  echo -e "${YELLOW}[info] stopping current zirva process...${NC}"
  kill -15 $ZIRVA_PID
fi

ARCH=$(uname -m)
if [ "$ARCH" = "i686" ] || [ "$ARCH" = "i386" ]; then
  ARCH_TYPE="i386"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "x86_64" ]; then
  ARCH_TYPE="$ARCH"
else
  echo -e "${RED}[err] unsupported architecture: $ARCH${NC}"
  exit 1
fi

PACKAGE_URL=$(curl -s $LATEST_RELEASE_URL | grep "browser_download_url.*$ARCH_TYPE" | cut -d '"' -f 4)

if [ -z "$PACKAGE_URL" ]; then
  echo -e "${RED}[err] no suitable package found for your architecture.${NC}"
  exit 1
fi

TEMP_DIR=$(mktemp -d)
PACKAGE_NAME=$(basename "$PACKAGE_URL")
echo -e "${YELLOW}[info] downloading the latest package...${NC}"
curl -L -o "$TEMP_DIR/$PACKAGE_NAME" "$PACKAGE_URL"

echo -e "${YELLOW}[info] installing the new version to $INSTALL_DIR...${NC}"
mv "$TEMP_DIR/$PACKAGE_NAME" "$INSTALL_DIR/zirva"
chmod +x "$INSTALL_DIR/zirva"

rm -rf "$TEMP_DIR"

echo -e "${GREEN}[success] update completed successfully! running new version...\n${NC}"
echo -e "${YELLOW}[info] starting the new version of zirva...${NC}"
clear

if [ -n "$ZIRVA_PORT" ]; then
  $INSTALL_DIR/zirva -p "$ZIRVA_PORT" &
else
  $INSTALL_DIR/zirva &
fi