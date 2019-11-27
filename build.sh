#!/bin/bash

# export DEBUG_LEVEL=30
# export DEBUG_DEFAULT_PACKAGE_LEVEL=30
# export DEBUG_DEFAULT_FUNCTION_LEVEL=30
# export DEBUG_PACKAGE_LEVEL_httphandler=30
# export DEBUG_FUNCTION_LEVEL_httphandler_CreateCourt=30

export GOPATH=${WORKSPACE}
export PROJECT_DIR=${WORKSPACE}/src/github.com/rsmaxwell/players-api

export PLAYERS_API_ROOTDIR=${PROJECT_DIR}/build/root
export DEBUG_DUMP_DIR=${PROJECT_DIR}/build/dumps

gradle clean generate build
