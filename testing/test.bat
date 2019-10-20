
rem set ENDPOINT=https://server.rsmaxwell.co.uk/players-api
set ENDPOINT=http://localhost:4201/players-api

set COMMAND=/register

(
	echo {
	echo     "userID": "007",
	echo     "first_name": "James",
	echo     "last_name": "Bond",
	echo     "email": "james@mi6.co.uk",
	echo     "password": "topsecret"
	echo }
) > data.json

curl -X POST %ENDPOINT%%COMMAND% ^
--header "Content-Type: application/json" ^
--data-binary @data.json