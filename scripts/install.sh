#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/zirva"
LATEST_RELEASE_URL="https://api.github.com/repos/zirvaorg/server/releases/latest"
UPDATE_SCRIPT_URL="https://raw.githubusercontent.com/zirvaorg/server/main/scripts/update.sh"

echo -e "${BLUE}... zirva server installer ...${NC}"

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[err] please run this script as root (using sudo).${NC}"
  exit 1
fi

if [ -f "$INSTALL_DIR/zirva" ]; then
  echo -e "${GREEN}[info] zirva is already installed at $INSTALL_DIR.${NC}"
  exit 1
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
  echo -e "${RED}[err] no suitable package found.${NC}"
  exit 1
fi

TEMP_DIR=$(mktemp -d)
PACKAGE_NAME=$(basename "$PACKAGE_URL")
echo -e "${YELLOW}[info] downloading the package...${NC}"
curl -L -o "$TEMP_DIR/$PACKAGE_NAME" "$PACKAGE_URL"

mkdir -p "$INSTALL_DIR"
echo -e "${YELLOW}[info] installing the package to $INSTALL_DIR...${NC}"
mv "$TEMP_DIR/$PACKAGE_NAME" "$INSTALL_DIR/zirva"
chmod +x "$INSTALL_DIR/zirva"

rm -rf "$TEMP_DIR"

echo -e "${YELLOW}[info] creating symlink for zirva in /usr/local/bin...${NC}"
ln -sf "$INSTALL_DIR/zirva" /usr/local/bin/zirva

echo -e "${YELLOW}[info] reloading your bash profile...${NC}"
source ~/.bashrc

echo -e "${YELLOW}[info] downloading update.sh...${NC}"
curl -L -o "$INSTALL_DIR/update.sh" "$UPDATE_SCRIPT_URL"
chmod +x "$INSTALL_DIR/update.sh"

echo -e "${YELLOW}[info] checking if update.sh is already in crontab...${NC}"
CRON_JOB="0 2 * * * $INSTALL_DIR/update.sh >> $INSTALL_DIR/zirva-update.log 2>&1"
(crontab -l 2>/dev/null | grep -q "$INSTALL_DIR/update.sh") || (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -

echo -e "${GREEN}[success] installation and cronjob setup completed successfully! running zirva...${NC}"
clear

$INSTALL_DIR/zirva