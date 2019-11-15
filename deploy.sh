#!/bin/bash

#**************************************************************
# Deploy 'players-api' to 'build'
#**************************************************************

currentDir=$(pwd)
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi

cd ${WORKSPACE}/bin
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi

#****************************************************************
#* Deploy the main executable and test utilities 
#****************************************************************
REPOSITORY_URL="https://server.rsmaxwell.co.uk/archiva/repository"
REPOSITORYID="build"
REPOSITORY="build"
GROUPID="com.rsmaxwell.players"
NAME="players-api"
PACKAGING="zip"
VERSION=${BUILD_ID}

GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)

URL=${REPOSITORY_URL}/${REPOSITORY}/
ARTIFACTID=${NAME}-${GOARCH}-${GOOS}
FILENAME=${ARTIFACTID}-${VERSION}.${PACKAGING}

CLASSIFIER_TEST="test"
FILENAME_TEST=${ARTIFACTID}-${VERSION}-${CLASSIFIER_TEST}.${PACKAGING}


zip ${FILENAME} players-api version.json
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi

zip ${FILENAME_TEST} generate-test-data
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi
mvn --batch-mode deploy:deploy-file \
	-Durl=${URL} \
	-DrepositoryId=${REPOSITORYID} \
	-Dfile=${FILENAME} \
	-DgroupId=${GROUPID} \
	-DartifactId=${ARTIFACTID} \
	-Dversion=${VERSION} \
	-Dpackaging=${PACKAGING} \
	-Dfiles=${FILENAME_TEST} \
	-Dclassifiers=${CLASSIFIER_TEST} \
	-Dtypes=zip \
    1>stdout.txt 2>stderr.txt

result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    echo "----[ stdout ]--------------------------"
    cat stdout.txt
    echo "----[ stderr ]--------------------------"
    cat stderr.txt
    echo "----------------------------------------"
    exit 1
fi

#****************************************************************
#* Restore the original directory
#****************************************************************
cd ${currentDir}
result=$?
if [ ! ${result} == 0 ]; then
    echo "Error: $0[${LINENO}]"
    echo "result: ${result}"
    exit 1
fi