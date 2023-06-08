#!/usr/bin/env bash
set -eu


DEST=~/Library/LaunchAgents

echo "copying plist to system directory"
cp jeffsnapshot.plist $DEST


#sudo chown root:wheel /Library/LaunchDaemons/jeffsnapshot.plist


echo "loading service"
launchctl load $DEST/jeffsnapshot.plist


echo "starting service"
launchctl start jeffsnapshot

echo "viewing in list"
launchctl list | grep jeffsnapshot
