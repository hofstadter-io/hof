#!/usr/bin/env bash
set -euo pipefail

# REGISTRY=host.docker.internal:5000
REGISTRY=localhost:5000
INSECURE=--insecure

# 1. Get Registry Images
# We pipe to tr -d '\r' to remove any invisible carriage returns
# We pipe to awk to trim leading/trailing whitespace
krane $INSECURE catalog $REGISTRY | tr -d '\r' | awk '{$1=$1};1' | sort > reg_sorted.txt

# 2. Get SQL Images
# -batch: prevents interactive prompts
# -noheader: ensures "id" isn't printed as the first line
# -csv: ensures clean output without padding
sqlite3 -batch -noheader -csv second-veg.db 'select distinct(eid) from environs' | tr -d '\r' | awk '{$1=$1};1' | sort > sql_sorted.txt

# 3. Compare
# -23: Suppress col 2 (unique to SQL) and col 3 (matches). 
# Result is col 1 (unique to Registry)
comm -23 reg_sorted.txt sql_sorted.txt > cmp.txt

# Output result to console so you can see it
echo "--- Removing Images in Registry but NOT in SQL ---"

for img in `cat cmp.txt`; do
  echo $img
  for tag in `krane $INSECURE ls $REGISTRY/$img`; do
    dig=$(krane $INSECURE digest $REGISTRY/$img:$tag)
    echo "  :$tag @$dig"
    # krane $INSECURE delete $REGISTRY/$img:$tag
    # krane $INSECURE delete $REGISTRY/$img@$dig
  done
done


# Clean up temp files
# rm reg_sorted.txt sql_sorted.txt

