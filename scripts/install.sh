#!/bin/bash

RED='\033[1;31m'
GREEN='\033[1;32m'
BLUE='\033[1;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/zirva"
LATEST_RELEASE_URL="https://api.github.com/repos/zirvaorg/server/releases/latest"
UPDATE_SCRIPT_URL="https://zirva.org/update.sh"
SERVER_PORT=9479

echo -e "${BLUE}... zirva server installer ...${NC}"

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[err] please run this script as root (using sudo).${NC}"
  exit 1
fi

if ! command -v curl &> /dev/null; then
  echo -e "${RED}[err] curl is not installed. please install curl and try again.${NC}"
  exit 1
fi

if ! command -v crontab &> /dev/null; then
  echo -e "${RED}[err] crontab is not installed. please install crontab and try again.${NC}"
  exit 1
fi

if [ -f "$INSTALL_DIR/zirva" ]; then
  echo -e "${GREEN}[info] zirva is already installed at $INSTALL_DIR.${NC}"
  exit 1
fi

if lsof -i:$SERVER_PORT &> /dev/null; then
  echo -e "${RED}[err] port $SERVER_PORT is already in use. please free the port and try again.${NC}"
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
echo -e "${BLUE}[info] downloading the package...${NC}"
curl -sL -o "$TEMP_DIR/$PACKAGE_NAME" "$PACKAGE_URL"

mkdir -p "$INSTALL_DIR"
echo -e "${BLUE}[info] installing the package to $INSTALL_DIR...${NC}"
mv "$TEMP_DIR/$PACKAGE_NAME" "$INSTALL_DIR/zirva"
chmod +x "$INSTALL_DIR/zirva"

rm -rf "$TEMP_DIR"

echo -e "${BLUE}[info] creating symlink for zirva in /usr/bin...${NC}"
ln -sf "$INSTALL_DIR/zirva" /usr/bin/zirva

echo -e "${BLUE}[info] reloading your bash profile...${NC}"
source ~/.bashrc

echo -e "${BLUE}[info] downloading update.sh...${NC}"
curl -sL -o "$INSTALL_DIR/update.sh" "$UPDATE_SCRIPT_URL"
chmod +x "$INSTALL_DIR/update.sh"

echo -e "${BLUE}[info] checking if update.sh is already in crontab...${NC}"
CRON_JOB="0 2 * * * $INSTALL_DIR/update.sh >> $INSTALL_DIR/update.log 2>&1"
(crontab -l 2>/dev/null | grep -q "$INSTALL_DIR/update.sh") || (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -

if command -v systemctl &> /dev/null; then
  echo -e "${BLUE}[info] creating systemd service file...${NC}"
  SERVICE_FILE="/etc/systemd/system/zirva.service"
  cat <<EOF > $SERVICE_FILE
[Unit]
Description=Zirva Server Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/zirva
Restart=always
User=root
WorkingDirectory=$INSTALL_DIR
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

  echo -e "${BLUE}[info] enabling zirva service to start on boot...${NC}"
  systemctl enable zirva.service
else
  echo -e "${YELLOW}[warn] systemd is not available. zirva will not start on boot.${NC}"
fi

echo -e "${GREEN}[ok] installation completed successfully! running zirva...${NC}"
$INSTALL_DIR/zirva