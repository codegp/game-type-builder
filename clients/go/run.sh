#!/bin/bash

set -e
cd /go/src/teamrunner
mkdir /source
mkdir -p /go/src/teamrunner/bot
echo "Running source manager..."
./sourcemanager
echo "Done running source manager"
ls /source
echo "YOUR"
for file in /source/*; do
  echo "HILLLE"
  echo $file
  FNAME=$(echo $file | cut -d'/' -f2)
  echo $FNAME
  mv "$file" "/go/src/teamrunner/bot/$FNAME.go"
done

echo "Running bot"
ls /go/src/teamrunner/bot
go build
./teamrunner
