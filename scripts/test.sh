#!/bin/bash
CURR_DIR=$(dirname "$0");
tmpfile="./tmp/ai" # $(mktemp);
bash "$CURR_DIR/create_release.sh" "$tmpfile"; 
bash "$tmpfile" "$@"; 
# rm "$tmpfile"