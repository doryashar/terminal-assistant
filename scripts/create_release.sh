#!/bin/bash
tmpfile=$1; printf "" > "$tmpfile";
for file in ./src/{common_funcs,config,push_to_stdin,conversations,agent,install,update,main}; do cat "$file" >> "$tmpfile"; echo >> "$tmpfile"; done; 