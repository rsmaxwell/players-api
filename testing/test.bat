
set ENDPOINT=https://server.rsmaxwell.co.uk/players-api
rem set ENDPOINT=http://localhost:4201/players-api

GOTO login

set COMMAND=/register

(
	echo {
	echo     "userID": "007",
	echo     "firstname": "James",
	echo     "lastname": "Bond",
	echo     "email": "james@mi6.co.uk",
	echo     "password": "topsecret"
	echo }
) > data.json

curl -k -X POST %ENDPOINT%%COMMAND% ^
--header "Content-Type: application/json" ^
--data-binary @data.json





:login
set COMMAND=/login
set USERID=007
set PASSWORD=topsecret

curl -k -X GET -u "%USERID%:%PASSWORD%" %ENDPOINT%%COMMAND% ^
--header "Content-Type: application/json" ^
--header "Accept: application/json"