#!/bin/bash

#redefine GOPATH
PCKT_PATH=$(echo $(pwd) | sed  -e 's/\/src//g')
GOPATH=$GOPATH:$PCKT_PATH

#set GOPATH
GOPATH=$GOPATH goapp serve

