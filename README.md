# Players
A Player manager - backend REST server

(see: https://thenewstack.io/make-a-restful-json-api-go/)

### Build
Get the dependancies:

```
make deps
```


### Install
The application data is stored in the "${HOME}/players-api" directory.


### Run

Given the following variables are set:
``` bash
USER=foo
PASSWORD=bar
ENDPOINT=localhost:4201

players-api
```

### List all people
``` bash
COMMAND="/person"

curl -X GET -u "${USER}:${PASSWORD}" ${ENDPOINT}${COMMAND} \
--header "Accept: application/json"
```

``` json
httpStatus: 200
response:   { "people":[1001,1002] }
```


### Add a new Person
``` bash
COMMAND="/person"

cat <<EOT > data.json
{
    "name": "xxx"
}
EOT

curl -X POST -u "${USER}:${PASSWORD}" ${ENDPOINT}${COMMAND} \
--header "Content-Type: application/json" \
--data-binary @data.json
```

``` json
httpStatus: 200
response:   { "message":"ok" }
```

### Delete a person
``` bash
COMMAND="/person"
ID=1002

curl -X DELETE -u "${USER}:${PASSWORD}" ${ENDPOINT}${COMMAND}/${ID}
```

``` json
httpStatus: 200
response:   {"message":"ok"}
```

### Get the details of a person
``` bash
COMMAND="/person"
ID=1002

curl -X GET -u "${USER}:${PASSWORD}" ${ENDPOINT}${COMMAND}/${ID}
```

``` json
httpStatus: 200
response:   { "person":{"name":"FRED"} }
```



