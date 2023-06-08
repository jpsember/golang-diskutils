#!/usr/bin/env bash
set -eu
home="/Users/home"
now=$(date "+%Y-%m-%d %H.%M.%S")
echo $now >> "$home/Desktop/TestFile.txt"
