#!/usr/bin/env bash

set -x # trace what gets executed. Useful for debugging
set -e # make your script exit when a command fails
set -u # exit when your script tries to use undeclared variables


source ./common.sh


# 0 - UPDATE VERSION IN DOCKER COMPOSE
VERSION=$(getVersion)

mkdir -p ../acceptance-tests/target && \
sed "s#VERSION#${VERSION}#g" ../deploy/acceptance-tests/docker-compose.yml > ../acceptance-tests/target/docker-compose.yml && \
cp -f ../deploy/acceptance-tests/rancher-compose.yml ../acceptance-tests/target


# 1 - LAUNCH STACK
#
docker-rancher-tools rancher-compose -p ${STACK_NAME} -r acceptance-tests/target/rancher-compose.yml -f acceptance-tests/target/docker-compose.yml up security-vault dummy vault-server -d
sleep 5
docker-rancher-tools rancher-compose -p ${STACK_NAME} -r acceptance-tests/target/rancher-compose.yml -f acceptance-tests/target/docker-compose.yml up vault-configurator
sleep 5
docker-rancher-tools rancher-compose -p ${STACK_NAME} -r acceptance-tests/target/rancher-compose.yml -f acceptance-tests/target/docker-compose.yml up security-vault-at -d


# 2 - GET LOGS
docker-rancher-api-cli -c "fl services,fr ${AT_DOCKER_NAME},fl instances,fr ${STACK_NAME}_${AT_DOCKER_NAME}_1,fa logs,rl" | tee ../acceptance-tests/target/rancher.log


# 3 - CLEAN STACK
docker-rancher-tools rancher-compose -p ${STACK_NAME} -r acceptance-tests/target/rancher-compose.yml -f acceptance-tests/target/docker-compose.yml rm --force


# 4 - PARSE LOGS AND RETURN EXIT CODE (see acceptance-tests Dockerfile entrypoint)
if ( tail -n 2 ../acceptance-tests/target/rancher.log | grep -ai success ); then
    echo "SUCCESS" && rm -f ../acceptance-tests/target/rancher.log && exit 0
else
    echo "FAILED" && rm -f ../acceptance-tests/target/rancher.log && exit 1
fi
