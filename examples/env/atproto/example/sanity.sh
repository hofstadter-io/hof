#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)

ENV_FILE=${ENV_FILE:-"$REPO_ROOT/examples/.env"}
source "${ENV_FILE}"

echo "Handle: $HANDLE"
echo "DID:    $DID"


# check the PLC
echo "Checking PLC"
curl -s "https://$PLC/$DID" | jq .

# check the PDS
echo "Checking PDS"
curl -s "https://$HANDLE/.well-known/atproto-did" && echo ""
curl -s "https://$PDS/xrpc/com.atproto.sync.getRepoStatus?did=$DID" | jq .
curl -s "https://$PDS/xrpc/com.atproto.repo.describeRepo?repo=$DID" | jq .
