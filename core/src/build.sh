#!/bin/bash

#redefine GOPATH
PCKT_PATH=$(echo $(pwd) | sed  -e 's/\/src//g')
GOPATH=$GOPATH:$PCKT_PATH

#set GOPATH
GOPATH=$GOPATH goapp build core/err \
    core/appl/auth \
    core/appl/chat \
    core/appl/invite \
    core/appl/member \
    core/infra/comm \
    core/infra/data

