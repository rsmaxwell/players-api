# Players
An example of a Player manager - backend REST server

see   https://thenewstack.io/make-a-restful-json-api-go/

### Build
Get the dependancies:

```
make deps
```


### Install
The application data is stored in the "players" directory under:
```
Windows:    /ProgramData
Linux:      /var/lib
MacOS:      /Library/Application Support
```
which needs to be created as root as follows (for example on linux) :

```
sudo mkdir /var/lib/players-test
sudo chmod 777 /var/lib/players-test

``` 


### Run

Given the following variables are set:
```
USER=user
PASSWORD=pass
HOST=localhost
```

List the IDs of all players
```
curl ${USER}:${PASSWORD}@${HOST}:8080/players

httpStatus: 200
response:   {"players":[1001,1002]}
```


Add a new Player
```
curl -X POST -d "{\"name\":\"xxx\"}" ${USER}:${PASSWORD}@${HOST}:8080/player

httpStatus: 200
response:   {"httpStatus":200,"message":"ok"}
```

Delete a player
```
ID=1002
curl -X DELETE ${USER}:${PASSWORD}@${HOST}:8080/player/${ID}

httpStatus: 200
response:   {"httpStatus":200,"message":"ok"}
```

Get the details of a player
```
curl ${USER}:${PASSWORD}@${HOST}:8080/player/${ID}

httpStatus: 200
response:   {"player":{"name":"FRED"}}
```



