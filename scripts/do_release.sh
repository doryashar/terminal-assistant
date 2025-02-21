
#!/bin/bash
# Get the version from the config file
eval $(head -1 src/config)
CURR_DIR=$(dirname "$0");

# Create the script
SCRIPT_PATH=$(mktemp)
bash "$CURR_DIR/create_release.sh" "$SCRIPT_PATH"; 

# Ensure the script file exists
if [ ! -f "$SCRIPT_PATH" ]; then
    echo "[ERROR] The file '$SCRIPT_PATH' does not exist!"
    exit 1
fi

# Ensure VERSION is set
if [ -z "$VERSION" ]; then
    echo "[ERROR] VERSION is not set!"
    exit 1
fi

# Commit all changes
git add .
git commit -m "Release version $VERSION"

# Push to the current branch
git push origin "$(git rev-parse --abbrev-ref HEAD)"

# Create and push the tag
git tag "$VERSION"
git push origin "$VERSION"

# Optional: Create a GitHub release and upload the script (requires GitHub CLI)
if command -v gh &>/dev/null; then
    gh release create "$VERSION" "$SCRIPT_PATH#ai" \
        --title "Release $VERSION" \
        --notes "Automated release for version $VERSION"
    echo "[INFO] GitHub release created with script upload."
else
    echo "[WARNING] GitHub CLI (gh) not found. Release not created."
fi

echo "[INFO] Version $VERSION committed, pushed, tagged, and script uploaded."