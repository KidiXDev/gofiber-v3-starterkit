#!/bin/bash

REPO_URL="https://github.com/KidiXDev/gofiber-v3-starterkit.git"
BRANCH="main"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
RESET='\033[0m'

clear

echo -e "${BLUE}"
echo "   GoFiber V3 Starter Pack Installer"
echo -e "${RESET}"

if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: git is not installed.${RESET}"
    exit 1
fi

echo -e "${GREEN}Enter project name (default: my-gofiber-app):${RESET}"
if [ -t 0 ]; then
    read -p "> " PROJECT_NAME
else
    read -p "> " PROJECT_NAME < /dev/tty
fi
PROJECT_NAME=${PROJECT_NAME:-my-gofiber-app}

if [ -d "$PROJECT_NAME" ]; then
    echo -e "${RED}Error: Directory '$PROJECT_NAME' already exists.${RESET}"
    exit 1
fi

echo ""
echo -e "${BLUE}Cloning repository into '$PROJECT_NAME'...${RESET}"
git clone --depth 1 "$REPO_URL" "$PROJECT_NAME"

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to clone repository.${RESET}"
    exit 1
fi

cd "$PROJECT_NAME" || exit

rm -rf .git
git init

rm -f install.sh install.ps1

echo -e "${GREEN}Repository cloned successfully!${RESET}"
echo ""

chmod +x rename-module.sh

./rename-module.sh

if [ $? -eq 0 ]; then
    rm -f rename-module.sh rename-module.bat
fi
