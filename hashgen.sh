#!/bin/sh

for folder in ./dist/*; do
  echo "$folder"
  if [ -d "$folder" ]; then
    for f in "$folder"/*; do
      shasum -a 256 "$f" >>"$f".sha256
    done
  fi
done