
#!/bin/bash
CURR_DIR=$(dirname "$0");
LATEST_RELEASE=$(git describe --tags --abbrev=0)

# Increment the version
MAJOR=$(echo $LATEST_RELEASE | cut -d. -f1)
MINOR=$(echo $LATEST_RELEASE | cut -d. -f2)
PATCH=$(echo $LATEST_RELEASE | cut -d. -f3)
NEW_PATCH=$((PATCH + 1))
VERSION="${MAJOR}.${MINOR}.${NEW_PATCH}"

# Update the config file
sed -i "1s/.*VERSION=.*/VERSION=\"$VERSION\"/" src/config

# Create the script
TMP_PATH=$(mktemp)
SCRIPT_PATH=$CURR_DIR/install
bash "$CURR_DIR/create_release.sh" "$TMP_PATH"; 
mv $TMP_PATH $SCRIPT_PATH
chmod +x $SCRIPT_PATH

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
    gh release create "$VERSION" "$SCRIPT_PATH#install" \
        --title "Release $VERSION" \
        --notes "Automated release for version $VERSION"
    echo "[INFO] GitHub release created with script upload."
else
    echo "[WARNING] GitHub CLI (gh) not found. Release not created."
fi
rm -rf $SCRIPT_PATH
echo "[INFO] Version $VERSION committed, pushed, tagged, and script uploaded."