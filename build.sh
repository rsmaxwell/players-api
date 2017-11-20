#!/bin/bash 

go build

go test -c -covermode=count -coverpkg ./...
./players.test -test.coverprofile coverage.cov
go tool cover -html=coverage.cov -o players.html

