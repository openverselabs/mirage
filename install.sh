#!/bin/bash

CYAN='\033[0;36m'
GREEN='\033[1;32m'
BLUE='\033[1;34m'
RED='\033[1;31m'
NC='\033[0m'

echo -e "${CYAN}"
echo "        Mirage v0.1.0 Installation Script        "
echo -e "${NC}"
echo -e "${GREEN}[+] Starting installation of Mirage...${NC}\n"

if ! command -v go &> /dev/null; then
    echo -e "${RED}[!] Error: Golang is not installed.${NC}"
    echo -e "Please install Go first (https://go.dev/doc/install) and try again."
    exit 1
fi

if ! command -v git &> /dev/null; then
    echo -e "${RED}[!] Error: Git is not installed.${NC}"
    exit 1
fi

echo -e "${BLUE}[*] Cloning repository...${NC}"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR" || exit
git clone -q https://github.com/openverselabs/mirage.git
cd mirage || exit

echo -e "${BLUE}[*] Building Mirage (optimizing for size)...${NC}"
if [ ! -f "go.mod" ]; then
    go mod init github.com/openverselabs/mirage &> /dev/null
fi
go build -ldflags="-s -w" -o mirage main.go

echo -e "${BLUE}[*] Moving binary to /usr/local/bin (may require sudo password)...${NC}"
sudo mv mirage /usr/local/bin/

cd / || exit
rm -rf "$TMP_DIR"

echo -e "\n${GREEN}[✔] Installation successful!${NC}"
echo -e "You can now run ${CYAN}mirage${NC} from anywhere in your terminal."
echo -e "Try it: ${CYAN}echo \"github.com\" | mirage${NC}"
