#!/bin/bash

export GOPATH=$(cd $(dirname $0); pwd)
go get -u github.com/mailgun/mailgun-go
go get -u github.com/gorilla/mux

echo "Be sure to run \". ./env.sh\" if you intend on hacking the source code."
echo "This will set your GOPATH environment variable accordingly."

