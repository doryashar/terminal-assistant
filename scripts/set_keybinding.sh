#!/bin/bash

# # Define keys & strings to be used
# key="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings"
# firstname="custom"

# # Get the current list of custom shortcuts
# array_str=$(dconf read "$key")

# # Remove the annotation hints if the array is empty
# current=$(echo "$array_str" | sed 's/^@as\s*//')

# # Find an available slot for the new keybinding
# n=0
# while [[ "$current" == *"$key/$firstname$n/"* ]]; do
#     ((n++))
# done

# # Create the new keybinding slot and add it to the list
# new="$key/$firstname$n/"
# current="$current $new"
# echo "current: $current new: $new"
# # Update the list of custom keybindings
# dconf write "$key" "$(echo "$current" | sed 's/ \+//g')"

# # Create the new shortcut with the provided arguments
# dconf write "$new" "name \"'$1'\""
# dconf write "$new" "command \"'$2'\""
# dconf write "$new" "binding \"'$3'\""

#Run example: ./set_keybinding.sh "My Script" "/path/to/my_script.sh" "<Super><Shift>s"

KEY1='org.gnome.settings-daemon.plugins.media-keys'
KEY2='.custom-keybinding'
KEY3='custom-keybindings'
list_of_customs=$(gsettings get $KEY1 $KEY3)

# Take the last result from the list:



gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ command 'script --command "flameshot gui" /dev/null'