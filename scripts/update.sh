#!/bin/bash

RED='\033[1;31m'
GREEN='\033[1;32m'
BLUE='\033[1;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/zirva"
LATEST_RELEASE_URL="https://api.github.com/repos/zirvaorg/server/releases/latest"
LATEST_VERSION=$(curl -s $LATEST_RELEASE_URL | grep '"tag_name":' | cut -d '"' -f 4)

echo -e "${BLUE}... zirva server updater ...${NC}"
echo -e "${BLUE}[info] checking for new version of zirva...${NC}"

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
  echo -e "${GREEN}[ok] you are already using the latest version: $CURRENT_VERSION${NC}"
  exit 1
else
  echo -e "${YELLOW}[warn] new version available: $LATEST_VERSION. updating...${NC}"
fi

echo -e "${BLUE}[info] checking if zirva is running on process...${NC}"
ZIRVA_PID=$(pgrep -f "zirva")
if [ -z "$ZIRVA_PID" ]; then
  echo -e "${GREEN}[ok] zirva is not running.${NC}"
else
  echo -e "${BLUE}[info] stopping current zirva process...${NC}"
  kill -15 $ZIRVA_PID
fi

echo -e "${BLUE}[info] checking if zirva is running on service...${NC}"
if command -v systemctl &> /dev/null; then
  echo -e "${BLUE}[info] stopping zirva service...${NC}"
  systemctl stop zirva.service
else
  echo -e "${YELLOW}[warn] systemctl is not available. skipping service stop.${NC}"
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
echo -e "${BLUE}[info] downloading the latest package...${NC}"
curl -sL -o "$TEMP_DIR/$PACKAGE_NAME" "$PACKAGE_URL"

echo -e "${BLUE}[info] installing the new version to $INSTALL_DIR...${NC}"
mv "$TEMP_DIR/$PACKAGE_NAME" "$INSTALL_DIR/zirva"
chmod +x "$INSTALL_DIR/zirva"

rm -rf "$TEMP_DIR"

echo -e "${GREEN}[ok] update completed successfully! starting new version...${NC}"
if command -v systemctl &> /dev/null; then
  systemctl start zirva.service
else
  nohup $INSTALL_DIR/zirva &>$INSTALL_DIR/zirva.log &
fi
