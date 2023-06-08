#!/usr/bin/env bash
set -eu


DEST=~/Library/LaunchAgents

echo "copying plist to system directory"
cp jeffsnapshot.plist $DEST

echo "loading service"
launchctl load $DEST/jeffsnapshot.plist


echo "starting service"
launchctl start jeffsnapshot

echo "viewing in list"
launchctl list | grep jeffsnapshot
