#!/bin/sh
echo "Running Main App Building"
if [ -f /opt/app/main.go ]; then { \
  echo "I do be building";
  cd /opt/app;
  go mod download;
  go mod verify;
  go build -v -o . ./...;
}; fi
echo "Main App Building Finished"