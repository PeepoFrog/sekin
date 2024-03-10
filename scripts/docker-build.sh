#!/bin/env bash

VERSION=$1

docker build -t sekai_rpc:$VERSION -f sekai.Dockerfile .
