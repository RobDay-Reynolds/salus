#!/bin/bash

bin=$(dirname $0)

cd $(dirname $bin)

echo '##############################'
echo 'FYI: ICMP test will fail due to docker for mac issue (https://github.com/docker/for-mac/issues/57)'
echo '##############################'

go get github.com/onsi/ginkgo/ginkgo
sysctl -w net.ipv4.ping_group_range="0 65535"


$bin/test-unit && $bin/test-integration