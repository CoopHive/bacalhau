#!/bin/bash
set -xeuo pipefail

# hack to build bacalhau when you have newer golang installed
docker run -v .:/go/bacalhau golang:1.20 bash -c "cd bacalhau; go build -buildvcs=false"
