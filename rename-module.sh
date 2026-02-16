#!/bin/bash

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
RESET='\033[0m'

clear

echo -e "${BLUE}"
echo "   ____       _____ _ _                 __   _______           "
echo "  / ___| ___ |  ___(_) |__   ___ _ __   \ \ / /___ /           "
echo " | |  _ / _ \| |_  | | '_ \ / _ \ '__|___\ \ / /|_ \           "
echo " | |_| | (_) |  _| | | |_) |  __/ | |_____\ V /___) |          "
echo "  \____|\___/|_|   |_|_.__/ \___|_|        \_/|____/           "
echo -e "${RESET}"
echo -e "${YELLOW}Welcome to the GoFiber V3 Starter Pack Wizard!${RESET}"
echo "----------------------------------------------------"

# Check if module name is provided as argument
if [ -z "$1" ]; then
    echo -e "${GREEN}Please enter your new module name (e.g., github.com/username/project):${RESET}"
    read -p "> " NEW_MODULE
else
    NEW_MODULE="$1"
fi

if [ -z "$NEW_MODULE" ]; then
    echo -e "${RED}Error: Module name cannot be empty.${RESET}"
    exit 1
fi

OLD_MODULE="gofiber-starterkit"

echo ""
echo -e "You are about to rename the module from:"
echo -e "${RED}$OLD_MODULE${RESET} -> ${GREEN}$NEW_MODULE${RESET}"
echo ""
read -p "Are you sure? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}Operation cancelled.${RESET}"
    exit 1
fi

echo ""
echo -e "${BLUE}Renaming module...${RESET}"

# Find and replace in all .go files
find . -type f -name "*.go" -exec sed -i "s|$OLD_MODULE|$NEW_MODULE|g" {} +

# Replace in go.mod
sed -i "s|module $OLD_MODULE|module $NEW_MODULE|g" go.mod

# Update the scripts themselves to prevent re-running with old name issues if they persist
sed -i "s|OLD_MODULE=\"$OLD_MODULE\"|OLD_MODULE=\"$NEW_MODULE\"|g" rename-module.sh
sed -i "s|set \"OLD_MODULE=$OLD_MODULE\"|set \"OLD_MODULE=$NEW_MODULE\"|g" rename-module.bat 2>/dev/null || true

echo -e "${GREEN}✔ Module renamed successfully!${RESET}"
echo ""
echo -e "${YELLOW}Next steps:${RESET}"
echo "1. Run 'go mod tidy' to update dependencies"
echo "2. Run 'go build' to verify the build"
echo "3. Copy .env.example to .env and configure your environment"
echo "4. Run 'go run .' to start the server"
echo ""
