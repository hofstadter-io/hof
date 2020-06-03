#!/bin/bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOT="$DIR/../.."

SONAR_TOKEN=${1:-$SONAR_TOKEN}
GIT_BRANCH=${2:-`git rev-parse --abbrev-ref HEAD | tr -d "\n"`}

if [ -z $SONAR_TOKEN ]; then
    echo please supply an SONAR_TOKEN env var or first argument to this script
    echo or add it as FERRUM_SONAR_TOKEN in your .profile or similar
    exit 1
fi

docker run \
  -e SONAR_TOKEN=$SONAR_TOKEN \
  -e SONAR_BRANCH=$GIT_BRANCH \
  -v "$PWD:/usr/src" \
  hof/sonar-scanner-cli
