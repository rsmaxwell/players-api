APP:=$(notdir $(CURDIR))
GOPATH:=$(shell go env GOPATH)
GOPACKAGES:=$(shell go list ./... | grep -E -v '/integration-test/|/vendor/')
GOHOSTOS:=$(shell go env GOHOSTOS)
GOHOSTARCH:=$(shell go env GOHOSTARCH)

GIT_COMMIT_SHA:="$(shell git rev-parse HEAD 2>/dev/null)"
GIT_REMOTE_URL:="$(shell git config --get remote.origin.url 2>/dev/null)"

# Jenkins vars. Set to `unknown` if the variable is not yet defined
BUILD_ID?=unknown
BUILD_NUMBER?=unknown

ifeq (${GOHOSTOS},windows)
BUILD_DATE:="$(shell cmd /C generate_timestamp)"
EXTENSION:=.exe
else
BUILD_DATE:="$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")"
EXTENSION:=
endif

.PHONY: all deps test coverage vet build
all: deps build

${GOPATH}/bin/golint:
	go get github.com/golang/lint/golint

${GOPATH}/bin/gocovmerge:
	go get github.com/wadey/gocovmerge

${GOPATH}/bin/gometalinter:
	go get -u github.com/alecthomas/gometalinter
	${GOPATH}/bin/gometalinter --install

deps: vendor
vendor: glide.yaml
	glide install --strip-vendor

build: deps
ifneq (${GOHOSTOS}-${GOHOSTARCH},darwin-386)
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o ${APP}-darwin-386 -ldflags "-s -w" -a -installsuffix cgo .
endif

ifneq (${GOHOSTOS}-${GOHOSTARCH},darwin-amd64)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${APP}-darwin-amd64 -ldflags "-s -w" -a -installsuffix cgo .
endif

ifneq (${GOHOSTOS}-${GOHOSTARCH},linux-386)
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ${APP}-linux-386 -ldflags "-s -w" -a -installsuffix cgo .
endif

ifneq (${GOHOSTOS}-${GOHOSTARCH},linux-amd64)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP}-linux-amd64 -ldflags "-s -w" -a -installsuffix cgo .
endif

ifneq (${GOHOSTOS}-${GOHOSTARCH},windows-amd64)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${APP}-windows-amd64.exe -ldflags "-s -w" -a -installsuffix cgo .
endif

	CGO_ENABLED=0 go build -o ${APP}-${GOHOSTOS}-${GOHOSTARCH}${EXTENSION} -ldflags "-s -w" -a -installsuffix cgo .


vet: ${GOPATH}/bin/gometalinter
	${GOPATH}/bin/gometalinter --disable-all --enable=gofmt --enable=golint --enable=vet --enable=vetshadow --enable=ineffassign --enable=goconst --tests  --vendor -e ./...

test: vet ${GOPATH}/bin/gocovmerge
ifneq (${GOHOSTOS}-${GOHOSTARCH},linux-386)
	go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}go test -race -test.v -test.timeout=120s -coverprofile={{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile {{.ImportPath}}{{end}}' $(GOPACKAGES) | xargs -I {} bash -c {}
else
	go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}go test -test.v -test.timeout=120s -coverprofile={{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile {{.ImportPath}}{{end}}' $(GOPACKAGES) | xargs -I {} bash -c {}
endif
	@${GOPATH}/bin/gocovmerge `ls *.coverprofile` > cover.out
	@rm -f *.coverprofile
	go tool cover -func=cover.out

coverage: cover.html
cover.html: cover.out
	go tool cover -html=cover.out -o=cover.html
