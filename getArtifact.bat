@echo off

set REPOSITORY=internal
set GROUPID=com.rsmaxwell.players-api
set ARTIFACTID=players-api-amd64-linux
set PACKAGING=zip
set VERSION=19
set URL=https://server.rsmaxwell.co.uk/archiva/repository/%REPOSITORY%

set GROUPID2=%GROUPID:.=\%
set artifact=%USERPROFILE%\.m2\repository\%GROUPID2%\%ARTIFACTID%\%VERSION%\%ARTIFACTID%-%VERSION%.%PACKAGING%
echo "artifact: %artifact%"

if not exist %artifact% (
    call mvn dependency:get -DgroupId=%GROUPID% -DartifactId=%ARTIFACTID% -Dversion=%VERSION% -Dpackaging=%PACKAGING% -DremoteRepositories=%URL%
)

dir %artifact%