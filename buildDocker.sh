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

# TEMP
rm -rf vendor/github.com/codegp
mkdir -p $REPO_ROOT/vendor/github.com/codegp
cp -r $GOPATH/src/github.com/codegp/cloud-persister $REPO_ROOT/vendor/github.com/codegp/cloud-persister
cp -r $GOPATH/src/github.com/codegp/env $REPO_ROOT/vendor/github.com/codegp/env
cp -r $GOPATH/src/github.com/codegp/job-client $REPO_ROOT/vendor/github.com/codegp/job-client
cp -r $GOPATH/src/github.com/codegp/game-object-types $REPO_ROOT/vendor/github.com/codegp/game-object-types

mkdir -p $REPO_ROOT/game-runner/vendor/github.com/codegp
cp -r $GOPATH/src/github.com/codegp/cloud-persister $REPO_ROOT/game-runner/vendor/github.com/codegp/cloud-persister
cp -r $GOPATH/src/github.com/codegp/env $REPO_ROOT/game-runner/vendor/github.com/codegp/env
cp -r $GOPATH/src/github.com/codegp/job-client $REPO_ROOT/game-runner/vendor/github.com/codegp/job-client
cp -r $GOPATH/src/github.com/codegp/game-object-types $REPO_ROOT/game-runner/vendor/github.com/codegp/game-object-types

docker build -t gcr.io/$PROJECT_ID/game-type-builder:$TAG .

if $PUSH; then
  gcloud docker push gcr.io/$PROJECT_ID/game-type-builder:$TAG
fi
