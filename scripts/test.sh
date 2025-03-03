#!/bin/bash
CURR_DIR=$(dirname "$0");
tmpfile="./tmp/ai" # $(mktemp);
bash "$CURR_DIR/create_release.sh" "$tmpfile"; 
AI_DEBUG=1 bash "$tmpfile" "$@"; 
# rm "$tmpfile"