#!/bin/bash

export GOPATH=${WORKSPACE}
cd src/github.com/rsmaxwell/players-api

export DEBUG_LEVEL=30
export DEBUG_DEFAULT_PACKAGE_LEVEL=30
export DEBUG_DEFAULT_FUNCTION_LEVEL=30

export DEBUG_DUMP_DIR=${WORKSPACE}/build/dumps

gradle clean generate build
