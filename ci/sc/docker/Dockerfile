FROM sonarsource/sonar-scanner-cli:4.4

USER root
RUN apk update && apk upgrade && apk add --no-cache bash tree
USER scanner-cli

COPY entrypoint.sh /usr/bin/entrypoint.sh
