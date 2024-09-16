#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}[info] checking for new version of zirva...${NC}"

INSTALL_DIR="/opt/zirva"

if [ -f "$INSTALL_DIR/zirva" ]; then
  CURRENT_VERSION=$("$INSTALL_DIR/zirva" --version)
else
  echo -e "${RED}[err] zirva is not installed.${NC}"
  exit 1
fi

LATEST_RELEASE_URL="https://api.github.com/repos/zirvaorg/backend/releases/latest"
LATEST_VERSION=$(curl -s $LATEST_RELEASE_URL | grep '"tag_name":' | cut -d '"' -f 4)

if [ "$CURRENT_VERSION" = "$LATEST_VERSION" ]; then
  echo -e "${GREEN}[success] you are already using the latest version: $CURRENT_VERSION${NC}"
  exit 0
else
  echo -e "${YELLOW}[info] new version available: $LATEST_VERSION. updating...${NC}"
fi

echo -e "${YELLOW}[info] stopping current zirva process...${NC}"
pkill -f "$INSTALL_DIR/zirva"

ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
  ARCH_TYPE="64bit"
elif [ "$ARCH" = "i686" ] || [ "$ARCH" = "i386" ]; then
  ARCH_TYPE="32bit"
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

echo -e "${YELLOW}[info] starting the new version of zirva...${NC}"
$INSTALL_DIR/zirva &

echo -e "${GREEN}[success] update completed successfully! running new version...${NC}"