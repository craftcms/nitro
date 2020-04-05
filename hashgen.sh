#!/bin/sh

for f in ./*; do shasum -a 256 $f > $f.sha256; done