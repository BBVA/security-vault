#!/usr/bin/env bash

set -x # trace what gets executed. Useful for debugging
set -e # make your script exit when a command fails
set -u # exit when your script tries to use undeclared variables

source ./common.sh

case "$1" in
    push-tag)
        #docker-sbt sbt push-tag
        cd ..
        git tag -f `cat version`
        git push --force origin --tags
      ;;

    publish-rancher-catalog)
        # 1 - Clone Liquidadores catalog
        git clone ${RANCHER_CATALOG_URI} ../target/catalog && \
        mkdir -p "../target/catalog/templates/${SERVICE_NAME}/${GO_PIPELINE_COUNTER}"

        # 2 - Include service files - new version
        VERSION=$(getVersion)

        #
        mkdir -p ../target/service

        sed "s#REPO#${REGISTRY_URL}#g; s#VERSION#${VERSION}#g" ../deploy/service/docker-compose.yml > ../target/service/docker-compose.yml
        sed "s#SERVICE_NAME#${SERVICE_NAME}#g; s#VERSION#${VERSION}#g" ../deploy/service/rancher-compose.yml > ../target/service/rancher-compose.yml

        cp -f ../target/service/docker-compose.yml ../target/catalog/templates/${SERVICE_NAME}/${GO_PIPELINE_COUNTER}/
        cp -f ../target/service/rancher-compose.yml ../target/catalog/templates/${SERVICE_NAME}/${GO_PIPELINE_COUNTER}/

        # 3 - Push changes
        (cd ../target/catalog ; git add -A && git diff-index --quiet HEAD || (git commit -a -m "feat (${VERSION}): New version" && git push origin master))
      ;;

    publish-api-docs)
        BUILD_UID=${BUILD_UID:-`id -u`}
        BUILD_GID=${BUILD_GID:-`getent group docker | cut -d: -f3`}

        docker run -i --rm -u ${BUILD_UID}:${BUILD_GID} -v $PWD/../documentation:/docs humangeo/aglio -i apiary.apib -o api-docs.html --theme-variables slate --verbose
      ;;

    *)
      echo -e "\n Option not recognized. Use $0 + 'push-tag', 'publish-rancher-catalog' or 'publish-api-docs'. \n"
      exit 1
      ;;
esac
