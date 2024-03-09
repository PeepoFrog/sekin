#!/bin/env bash

VERSION=$1
docker run --name sekai_rpc --rm -d -p 8080:8080 -v /home/0xddd/Code/github.com/mrlutik/sekin/sekai:/sekai sekai_rpc:$VERSION
