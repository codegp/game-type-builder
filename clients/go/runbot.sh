#!/bin/bash

set -e
cd /go/src/botrunner
echo "Running source manager..."
./sourcemanager
echo "Done running source manager"
ls /source
for file in /source/*; do
  FNAME=$(echo $file | cut -d'/' -f2)
  mv "$file" "/go/src/botrunner/bot/$FNAME.go"
done

echo "Running bot"
ls /go/src/botrunner/bot
go build
./botrunner
