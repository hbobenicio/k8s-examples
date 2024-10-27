#!/bin/bash
set -eu -o pipefail

set -x
docker build . -f ./docker/dev/Dockerfile     -t my-operator:0.0.0-dev
docker build . -f ./docker/release/Dockerfile -t my-operator:0.0.0
{ set +x; } 2>/dev/null
