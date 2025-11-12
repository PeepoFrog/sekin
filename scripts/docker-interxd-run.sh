#!/bin/env bash

VERSION=$1
docker run --name interx_rpc --rm -d -p 11000:11000 -p 8080:8080 -v $(pwd)/interx:/interx interx_rpc:$VERSION
