#!/bin/bash

export GOPATH=${WORKSPACE}

# export DEBUG_LEVEL=30
# export DEBUG_DEFAULT_PACKAGE_LEVEL=30
# export DEBUG_DEFAULT_FUNCTION_LEVEL=30
# export DEBUG_PACKAGE_LEVEL_httphandler=30
# export DEBUG_FUNCTION_LEVEL_httphandler_CreateCourt=30

export DEBUG_DUMP_DIR=${WORKSPACE}/src/github.com/rsmaxwell/players-api/build/dumps

gradle clean generate build
