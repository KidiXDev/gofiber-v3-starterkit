#!/bin/bash

set -e

if [ -z "$1" ]; then
    echo "Usage: ./rename-module.sh <new-module-name>"
    echo "Example: ./rename-module.sh github.com/username/my-project"
    exit 1
fi

OLD_MODULE="gofiber-starterkit"
NEW_MODULE="$1"

echo "Renaming module from '$OLD_MODULE' to '$NEW_MODULE'..."

find . -type f -name "*.go" -exec sed -i "s|$OLD_MODULE|$NEW_MODULE|g" {} +

sed -i "s|module $OLD_MODULE|module $NEW_MODULE|g" go.mod

sed -i "s|OLD_MODULE=\"$OLD_MODULE\"|OLD_MODULE=\"$NEW_MODULE\"|g" rename-module.sh
sed -i "s|set \"OLD_MODULE=$OLD_MODULE\"|set \"OLD_MODULE=$NEW_MODULE\"|g" rename-module.bat 2>/dev/null || true

echo "Module renamed successfully!"
echo ""
echo "Next steps:"
echo "1. Run 'go mod tidy' to update dependencies"
echo "2. Run 'go build' to verify the build"
echo "3. Copy .env.example to .env and configure your environment"
echo "4. Run the SQL migrations in migrations/ folder"
echo "5. Run 'go run .' to start the server"
