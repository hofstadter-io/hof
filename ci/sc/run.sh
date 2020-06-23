#!/bin/bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOT="$DIR/../.."

USAGE="
./ci/sonar_cloud/run.sh [-ih] [token] [branch]
  -h  help
  -i  interactive

  [token] as arg or SONAR_TOKEN
  [branch] as arg or SONAR_BRANCH, defaults to current git branch
"

INTERACTIVE="-t"
while getopts "i" OPTION; do
  case "${OPTION}" in
    i)
      INTERACTIVE="-it"
      ;;

    h)
      echo $USAGE
      exit 0
      ;;

    *)
      echo $USAGE
      exit 1
      ;;
  esac
done
shift $((OPTIND - 1))

# This allows supplying value with ENV var, as arg, and a help message in one line
SONAR_TOKEN=${1:-${SONAR_TOKEN:?"missing SonarCloud API token, supply as next arg or SONAR_TOKEN"}}
GIT_BRANCH=${2:-`git rev-parse --abbrev-ref HEAD | tr -d "\n"`}

# Extract the pull request id
PR_FLAG=""
if [[ "${CIRCLE_PULL_REQUEST:-}" != "" ]]; then
    PR_ID=${CIRCLE_PULL_REQUEST##*/}
    PR_FLAG="-e SONAR_PULL_REQUEST=$PR_ID"
fi

echo "================= env ======================="
env
echo "=================     ======================="

# Run the scanning container
echo docker run \
  $INTERACTIVE \
  --user $(id -u):$(id -g) \
  -e SONAR_TOKEN=$SONAR_TOKEN \
  -e SONAR_BRANCH=$GIT_BRANCH \
  $PR_FLAG \
  -e SONAR_CONFIG_FILE="./ci/sc/sonar-project.properties" \
  -v "$ROOT:/usr/src" \
  hof/sonar-scanner-cli

docker run \
  $INTERACTIVE \
  --user $(id -u):$(id -g) \
  -e SONAR_TOKEN=$SONAR_TOKEN \
  -e SONAR_BRANCH=$GIT_BRANCH \
  $PR_FLAG \
  -e SONAR_CONFIG_FILE="./ci/sc/sonar-project.properties" \
  -v "$ROOT:/usr/src" \
  hof/sonar-scanner-cli
