#!/usr/bin/env bash
set -eu

clear

echo "Making"
(cd ../..; mk.sh)

echo "Installing"
install.sh

echo "Watching logs"
tail -f /var/log/system.log
