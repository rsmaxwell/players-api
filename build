#!/bin/bash

export GOPATH=${WORKSPACE}

#******************************************************************
# Clean the 'bin' directory
#******************************************************************
rm -rf ${WORKSPACE}/bin
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi

mkdir -p ${WORKSPACE}/bin
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi

#******************************************************************
# Write out build info
#******************************************************************
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

find . -name "version.go" | while read versionfile; do

    echo "Replacing tags in ${versionfile}"

    sed -i "s@<BUILD_ID>@${BUILD_ID}@g"     ${versionfile}
    sed -i "s@<TIMESTAMP>@${TIMESTAMP}@g"   ${versionfile}
    sed -i "s@<GIT_COMMIT>@${GIT_COMMIT}@g" ${versionfile}
    sed -i "s@<GIT_BRANCH>@${GIT_BRANCH}@g" ${versionfile}
    sed -i "s@<GIT_URL>@${GIT_URL}@g"       ${versionfile}
done

x=$(echo "{}"   | jq --arg foo "${BUILD_ID}"   '. + {version: $foo}')
x=$(echo "${x}" | jq --arg foo "${TIMESTAMP}"  '. + {buildDate: $foo}')
x=$(echo "${x}" | jq --arg foo "${GIT_COMMIT}" '. + {gitCommit: $foo}')
x=$(echo "${x}" | jq --arg foo "${GIT_BRANCH}" '. + {gitBranch: $foo}')
x=$(echo "${x}" | jq --arg foo "${GIT_URL}"    '. + {gitURL: $foo}')

echo "${x}" > ${WORKSPACE}/bin/version.json

echo "---[ ${WORKSPACE}/bin/version.json ]------------------"
jq . ${WORKSPACE}/bin/version.json

#******************************************************************
# Start the build 
#******************************************************************
go get -d ./...
go install ./...