#!/bin/bash

cd $(cd "$(dirname "$0")"; pwd)

web_cmd="yarn --cwd web/"

if [ ${param_npm_repo} ]; then
  web_cmd+=" --registry ${param_npm_repo}"
fi

if [ ${param_yarn_cache} ]; then
  web_cmd+=" --cache-folder ${param_yarn_cache}"
fi

set -x

${web_cmd} install --no-lockfile --update-checksums --pure-lockfile && \
${web_cmd} build --ignore-engines --pure-lockfile && \
$(cd api/;go build -x -v -o build/ main.go)