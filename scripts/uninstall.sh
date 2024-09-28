#!/bin/bash

RED='\033[1;31m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/zirva"
SERVICE_FILE="/etc/systemd/system/zirva.service"

echo -e "${YELLOW}... zirva server uninstaller ...${NC}"

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[err] please run this script as root (using sudo).${NC}"
  exit 1
fi

if [ ! -f "$INSTALL_DIR/zirva" ]; then
  echo -e "${RED}[err] zirva is not installed at $INSTALL_DIR.${NC}"
  exit 1
fi

echo -e "${YELLOW}[info] stopping zirva if it's running...${NC}"
ZIRVA_PID=$(pgrep -f "zirva")
if [ -n "$ZIRVA_PID" ]; then
  echo -e "${YELLOW}[info] stopping zirva process with PID $ZIRVA_PID...${NC}"
  kill -15 "$ZIRVA_PID"
else
  echo -e "${GREEN}[success] zirva is not currently running.${NC}"
fi

if command -v systemctl &> /dev/null && [ -f "$SERVICE_FILE" ]; then
  echo -e "${YELLOW}[info] disabling and removing zirva systemd service...${NC}"
  systemctl stop zirva.service
  systemctl disable zirva.service
  rm -f "$SERVICE_FILE"
  systemctl daemon-reload
else
  echo -e "${YELLOW}[info] systemd service file not found or systemctl not available.${NC}"
fi

echo -e "${YELLOW}[info] removing zirva files from $INSTALL_DIR...${NC}"
rm -rf "$INSTALL_DIR"

echo -e "${YELLOW}[info] removing zirva symlink from /usr/bin...${NC}"
rm -f /usr/bin/zirva

echo -e "${YELLOW}[info] removing zirva update.sh from crontab...${NC}"
crontab -l 2>/dev/null | grep -v "$INSTALL_DIR/update.sh" | crontab -

echo -e "${GREEN}[success] zirva has been successfully uninstalled!${NC}"
