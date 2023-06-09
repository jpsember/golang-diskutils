#!/usr/bin/env bash
set -eu

clear

echo "Deleting old stderr, stdout"
rm -f stderr.txt stdout.txt

echo "Pulling"
(cd ../..;repull.sh)

echo "Making"
(cd ../..; mk.sh)

echo "Installing"
install.sh

echo "Watching logs"
tail -f /var/log/system.log
