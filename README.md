# Players
An example of a Player manager - backend REST server

see   https://thenewstack.io/make-a-restful-json-api-go/

### Build
Get the dependancies:

```
make deps
```


### Install
The application data is stored in the "players" directory under the "HOME" directory


### Run

Given the following variables are set:
```
USER=user
PASSWORD=pass
HOST=localhost
```

List the IDs of all players
```
curl ${USER}:${PASSWORD}@${HOST}:8080/people

httpStatus: 200
response:   {"people":[1001,1002]}
```


Add a new Person
```
curl -X POST -d "{\"name\":\"xxx\"}" ${USER}:${PASSWORD}@${HOST}:8080/person

httpStatus: 200
response:   {"httpStatus":200,"message":"ok"}
```

Delete a person
```
ID=1002
curl -X DELETE ${USER}:${PASSWORD}@${HOST}:8080/person/${ID}

httpStatus: 200
response:   {"httpStatus":200,"message":"ok"}
```

Get the details of a person
```
curl ${USER}:${PASSWORD}@${HOST}:8080/person/${ID}

httpStatus: 200
response:   {"person":{"name":"FRED"}}
```



