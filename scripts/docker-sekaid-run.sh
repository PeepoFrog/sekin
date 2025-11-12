#!/bin/env bash

VERSION=$1
docker run --name sekai_rpc --rm -d -p 8181:8080 -p 9090:9090 -p 26657:26657 -p 26656:26656 -v $(pwd)/sekai:/sekai sekai_rpc:$VERSION
