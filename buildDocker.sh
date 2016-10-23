#!/bin/bash

set -e

PROJECT_ID=$(gcloud config list project --format "value(core.project)" 2> /dev/null)
PUSH=false
TAG="latest"

while [ "$1" != "" ]; do
    PARAM=`echo $1 | awk -F= '{print $1}'`
    VALUE=`echo $1 | awk -F= '{print $2}'`
    case $PARAM in
        -d |--push)
            PUSH=$VALUE
            ;;
        -t | --tag)
            TAG=$VALUE
            ;;
        *)
            echo "ERROR: unknown parameter \"$PARAM\""
            exit 1
            ;;
    esac
    shift
done

REPO_ROOT=$GOPATH/src/github.com/codegp/game-type-builder

function cleanup {
  rm -rf $REPO_ROOT/game-runner
}
trap cleanup EXIT

cp -r $GOPATH/src/github.com/codegp/game-runner $REPO_ROOT/game-runner

docker build -t gcr.io/$PROJECT_ID/game-type-builder:$TAG .

if $PUSH; then
  gcloud docker push gcr.io/$PROJECT_ID/game-type-builder:$TAG
fi
