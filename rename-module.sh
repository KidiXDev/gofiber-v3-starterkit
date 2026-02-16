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


OS="$(uname)"
if [ "$OS" = "Darwin" ]; then
    SED_CMD="sed -i ''"
else
    SED_CMD="sed -i"
fi

if [ -f "go.mod" ]; then
    OLD_MODULE=$(grep "^module" go.mod | awk '{print $2}')
else
    echo -e "${RED}Error: go.mod not found. Cannot determine old module name.${RESET}"
    exit 1
fi

if [ -z "$OLD_MODULE" ]; then
    echo -e "${RED}Error: Could not determine module name from go.mod.${RESET}"
    exit 1
fi

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

run_sed() {
    local pattern="$1"
    local file="$2"
    if [ "$OS" = "Darwin" ]; then
        sed -i '' "$pattern" "$file"
    else
        sed -i "$pattern" "$file"
    fi
}

find . -type f \( -name "*.go" -o -name "go.mod" -o -name "*.md" -o -name "*.sh" -o -name "*.bat" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" \) -not -path "*/.*" -print0 | while IFS= read -r -d '' file; do
    if [[ "$file" == "./rename-module.sh" ]]; then
        continue
    fi
    run_sed "s|$OLD_MODULE|$NEW_MODULE|g" "$file"
done

run_sed "s|OLD_MODULE=\"$OLD_MODULE\"|OLD_MODULE=\"$NEW_MODULE\"|g" rename-module.sh
if [ -f "rename-module.bat" ]; then
    run_sed "s|$OLD_MODULE|$NEW_MODULE|g" rename-module.bat
fi

echo -e "${GREEN}✔ Module renamed successfully!${RESET}"
echo ""
echo -e "${YELLOW}Next steps:${RESET}"
echo "1. Run 'go mod tidy' to update dependencies"
echo "2. Run 'go build' to verify the build"
echo "3. Copy .env.example to .env and configure your environment"
echo "4. Run 'go run .' to start the server"
echo ""
