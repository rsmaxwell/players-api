@echo off

set WORKSPACE=C:\Users\Richard\go\src\github.com\rsmaxwell\players-api
set BUILD_ID=1001
set GIT_COMMIT=1002
set GIT_BRANCH=1003
set GIT_URL=1004

rem call gradle updateVersionGO

rem dir .\internal\basic\version\version.go
rem echo ---[ .\internal\basic\version\version.go ]-----------------------------
rem type .\internal\basic\version\version.go
rem echo ------------------------------------------------------------------------

call gradle createBinDir
