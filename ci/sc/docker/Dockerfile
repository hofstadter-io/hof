FROM sonarsource/sonar-scanner-cli

USER root
RUN apt-get update \
    && apt-get install -y --no-install-recommends tree \
    && rm -rf /var/lib/apt/lists/*
USER scanner-cli

COPY entrypoint.sh /usr/bin/entrypoint.sh
