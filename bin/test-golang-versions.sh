#!/bin/bash

set -xe

for GOLANG_VERSION in 1.6 1.7 1.8; do
    docker run --privileged -t -i -v $GOPATH:/gopath golang:$GOLANG_VERSION /gopath/src/github.com/monkeyherder/salus/bin/docker-test.sh
done