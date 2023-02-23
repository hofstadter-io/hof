#!/usr/bin/env bash
set -euo pipefail

IN=$1
OUTDIR=${2:-"."}

filename=$(basename -- "$IN")
extension="${filename##*.}"
filename="${filename%.*}"

echo "$filename $extension -> $OUTNAME-*.png"

echo convert -verbose -size 16x16 $IN $OUTDIR/

favicon-16x16.png
favicon-32x32.png
android-36x36.png
android-48x48.png
android-72x72.png
android-96x96.png
android-144x144.png
android-192x192.png
