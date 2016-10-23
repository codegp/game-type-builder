#!/bin/bash
set -e
usage()
{
cat << EOF
usage: $0 options

This script run the test1 or test2 over a machine.

OPTIONS:
   -h      Show this message
   -v      Docker tag
   -g      Game Type
EOF
}

while getopts “hg:v:” OPTION
do
     case $OPTION in
         h)
             usage
             exit 1
             ;;
         g)
             GAMETYPE=$OPTARG
             ;;
         v)
             TAG=$OPTARG
             ;;
         ?)
             usage
             exit
             ;;
     esac
done

# todo add guard to make sure we got flags
rm -rf tmp/
mkdir tmp
thrift -r -o tmp/ --gen py ../../thrift/api.thrift
thrift -r -o tmp/ --gen py ../../thrift/turnInformer.thrift
thrift -r -o tmp/ --gen py ../../thrift/$GAMETYPE.thrift
cp Dockerfile tmp/Dockerfile
cd ../../sourcemanager
CGO_ENABLED=0 GOOS=linux go build -o ../clients/python/tmp/sourcemanager
cd ../clients/python/tmp
docker build -t codegp/madvillainsminions:$GAMETYPE-pyclient-$TAG .
cd ..
rm -rf tmp/
