#!/bin/bash

export GOPATH=${WORKSPACE}
cd src/github.com/rsmaxwell/players-api

gradle zip publish
