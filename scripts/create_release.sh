#!/bin/bash
tmpfile=$1; printf "" > "$tmpfile";
for file in ./src/{common_funcs,config,agent,install,update,main}; do cat "$file" >> "$tmpfile"; echo >> "$tmpfile"; done; 
