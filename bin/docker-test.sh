#!/bin/bash

bin=$(dirname $0)

cd $(dirname $bin)

go get github.com/onsi/ginkgo/ginkgo
sysctl -w net.ipv4.ping_group_range="0 65535"


$bin/test-unit && $bin/test-integration