#!/bin/bash

#redefine GOPATH
PCKT_PATH=$(echo $(pwd) | sed  -e 's/\/src//g')
GOPATH=$GOPATH:$PCKT_PATH

#rename config
mv config.go config.go.test
mv app.yaml app.yaml.test
cp profile/config.go.prod config.go
cp profile/app.yaml.prod app.yaml

#set GOPATH
GOPATH=$GOPATH goapp deploy

#rename config
rm config.go
rm app.yaml
mv config.go.test config.go
mv app.yaml.test app.yaml
