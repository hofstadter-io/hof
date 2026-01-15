#!/usr/bin/env bash
set -euo pipefail

IN=$1
OUTDIR=${2:-"."}

convert -verbose -size 32x32 $IN $OUTDIR/favicon.ico
convert -verbose -size 16x16 $IN $OUTDIR/favicon-16x16.png
convert -verbose -size 32x32 $IN $OUTDIR/favicon-32x32.png
convert -verbose -size 36x36 $IN $OUTDIR/android-36x36.png
convert -verbose -size 48x48 $IN $OUTDIR/android-48x48.png
convert -verbose -size 72x72 $IN $OUTDIR/android-72x72.png
convert -verbose -size 96x96 $IN $OUTDIR/android-96x96.png
convert -verbose -size 144x144 $IN $OUTDIR/android-144x144.png
convert -verbose -size 192x192 $IN $OUTDIR/android-192x192.png
convert -verbose -size 180x180 $IN $OUTDIR/apple-touch-icon-180x180.png

