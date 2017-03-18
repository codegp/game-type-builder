#!/bin/bash

set -e

docker version
echo "IS_LOCAL $IS_LOCAL"

REPO_ROOT="$GOPATH/src/github.com/codegp/game-type-builder"
GAME_RUNNER_PATH_FROM_GOPATH_SRC="github.com/codegp/game-runner"
GAME_RUNNER_PATH="$GOPATH/src/$GAME_RUNNER_PATH_FROM_GOPATH_SRC"

echo "Building game type with params:"
echo "GCLOUD_PROJECT_ID: $GCLOUD_PROJECT_ID"
echo "GAME_TYPE_ID: $GAME_TYPE_ID"
echo "REPO_ROOT: $REPO_ROOT"
echo "GAME_RUNNER_PATH: $GAME_RUNNER_PATH"
echo "GAME_RUNNER_PATH_FROM_GOPATH_SRC: $GAME_RUNNER_PATH_FROM_GOPATH_SRC"

echo "Generating thrift files..."
docker ps
echo "afterpds"
ls /
echo "fo root"
cd /localstore
ls
echo "mo nothing"
cd $REPO_ROOT/thriftgenerator
mkdir $REPO_ROOT/thrift
cp $GAME_RUNNER_PATH/thrift/gameObjects.thrift $REPO_ROOT/thrift/gameObjects.thrift
go build
./thriftgenerator
echo "Done generating thrift files..."

echo "Generating thrift html ..."
cd $REPO_ROOT
thrift -o $REPO_ROOT/docsreporter --gen html:standalone thrift/ids.thrift
thrift -o $REPO_ROOT/docsreporter --gen html:standalone thrift/api.thrift
echo "Done generating thrift html..."

echo "Reporting html ..."
cd $REPO_ROOT/docsreporter
go build
./docsreporter
echo "Done reporting html..."

echo "Generating thrift go code ..."
cd $REPO_ROOT
thrift -r -out $GAME_RUNNER_PATH --gen go:package_prefix=$GAME_RUNNER_PATH_FROM_GOPATH_SRC/ thrift/api.thrift
echo "Done generating thrift go code..."

echo "Building gamerunner binary..."
cd $GAME_RUNNER_PATH/gamerunner
CGO_ENABLED=0 go build
echo "Done building gamerunner binary."

echo "Building gamerunner docker image..."
cd $GAME_RUNNER_PATH
docker build -t gcr.io/$GCLOUD_PROJECT_ID/game-runner-$GAME_TYPE_ID .
echo "Done building gamerunner docker image."

echo "Building sourcemanager..."
cd $REPO_ROOT/sourcemanager
CGO_ENABLED=0 go build
echo "Done building sourcemanager"

echo "Building client docker image..."
cd $REPO_ROOT
BOT_PATH="teamrunner"
thrift -r -out $REPO_ROOT/clients/go --gen go:package_prefix=$BOT_PATH/ thrift/ids.thrift
thrift -r -out $REPO_ROOT/clients/go --gen go:package_prefix=$BOT_PATH/ thrift/api.thrift
thrift -r -out $REPO_ROOT/clients/go --gen go:package_prefix=$BOT_PATH/ $GAME_RUNNER_PATH/thrift/turnInformer.thrift
thrift -r -out $REPO_ROOT/clients/go --gen go:package_prefix=$BOT_PATH/ $GAME_RUNNER_PATH/thrift/gameObjects.thrift
docker build -t gcr.io/$GCLOUD_PROJECT_ID/team-runner-$GAME_TYPE_ID-go -f clients/go/Dockerfile .
echo "Done building client docker image."
