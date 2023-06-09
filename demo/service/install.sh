#!/usr/bin/env bash
set -eu

plist_path="wtf.plist"
plist_filename=$(basename "$plist_path")
install_path="$HOME/Library/LaunchAgents"

echo "installing launchctl plist: $plist_path --> $install_path"

cp -f "$plist_path" "$install_path"

#sudo chown root "$install_path"
#sudo chmod 644 "$install_path"

echo "This may complain if it isn't already loaded..."
launchctl unload "$install_path"

echo "...now attempting to load it:"
launchctl load "$install_path"

echo "to check if it's running, run this command:  launchctl list | grep wtf"
echo "to uninstall, run this command:  launchctl unload \"$install_path\""
