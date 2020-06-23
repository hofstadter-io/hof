#!/usr/bin/env bash
set -euo pipefail

echo "starting scan..."

declare -a args

add_env_var_as_env_prop() {
  if [ "$1" ]; then
    args+=("-D$2=$1")
  fi
}

add_env_var_as_env_prop "${SONAR_LOGIN:-}" "sonar.login"
add_env_var_as_env_prop "${SONAR_PASSWORD:-}" "sonar.password"
add_env_var_as_env_prop "${SONAR_USER_HOME:-}" "sonar.userHome"
add_env_var_as_env_prop "${SONAR_PROJECT_BASE_DIR:-}" "sonar.projectBaseDir"

# if not empty, setup pr scan
if [[ "${SONAR_PULL_REQUEST:-}" != "" ]]; then
    add_env_var_as_env_prop "${SONAR_BRANCH:-}" "sonar.pullrequest.branch"
    add_env_var_as_env_prop "${SONAR_PULL_REQUEST:-}" "sonar.pullrequest.key"
else
    # otherwise branch scan
    add_env_var_as_env_prop "${SONAR_BRANCH:-}" "sonar.branch.name"
fi

PROJECT_BASE_DIR="$PWD"
if [ "${SONAR_PROJECT_BASE_DIR:-}" ]; then
  PROJECT_BASE_DIR="${SONAR_PROJECT_BASE_DIR}"
fi

# The grep inside the pipes is probably hiding errors
GO_TESTS="$(tree -f -i $PROJECT_BASE_DIR/{cmd,gen,lib,script} | grep "tests.json" | tr "\n" ",")sonar-reports/dummy.txt"
add_env_var_as_env_prop "${GO_TESTS:-}" "sonar.go.tests.reportPaths"

SONAR_CONFIG_FILE=${SONAR_CONFIG_FILE:-sonar-project.properties}
add_env_var_as_env_prop "${SONAR_CONFIG_FILE:-}" "project.settings"

echo "------- sonar args ------------"
echo "${args[@]}"
echo "------- sonar config ------------"
cat $SONAR_CONFIG_FILE
echo "---------------------------------"


sonar-scanner "${args[@]}"
