#!/bin/env bash

VERSION=$1

docker build -t interx_rpc:$VERSION -f interx.Dockerfile .
