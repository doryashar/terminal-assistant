#!/bin/bash
tmpfile=$1; printf "" > "$tmpfile";
for file in ./src/{common_funcs,config,terminal_buffer,push_to_stdin,conversations,agent,install,update,main,post_main}; do cat "$file" >> "$tmpfile"; echo >> "$tmpfile"; done; 