APP:=$(notdir $(CURDIR))
GOPATH:=$(shell go env GOPATH)
GOPACKAGES:=$(shell go list ./... | grep -E -v '/integration-test/|/vendor/')
GOHOSTOS:=$(shell go env GOHOSTOS)
GOHOSTARCH:=$(shell go env GOHOSTARCH)

GIT_COMMIT_SHA:="$(shell git rev-parse HEAD 2>/dev/null)"
GIT_REMOTE_URL:="$(shell git config --get remote.origin.url 2>/dev/null)"
BUILD_DATE:="$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")"
# Jenkins vars. Set to `unknown` if the variable is not yet defined
BUILD_ID?=unknown
BUILD_NUMBER?=unknown

.PHONY: all deps test coverage vet build
all: deps build

${GOPATH}/bin/golint:
	go get github.com/golang/lint/golint

${GOPATH}/bin/gocovmerge:
	go get github.com/wadey/gocovmerge

${GOPATH}/bin/gometalinter:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

deps: vendor
vendor: glide.yaml
	glide install --strip-vendor

build: deps
ifneq (${GOHOSTOS}-${GOHOSTARCH},linux-386)
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ${APP}-linux-386 -ldflags "-s -w" -a -installsuffix cgo .
endif
	CGO_ENABLED=0 go build -o ${APP}-${GOHOSTOS}-${GOHOSTARCH} -ldflags "-s -w" -a -installsuffix cgo .

vet: ${GOPATH}/bin/gometalinter
	${GOPATH}/bin/gometalinter --disable-all --enable=gofmt --enable=golint --enable=vet --enable=vetshadow --enable=ineffassign --enable=goconst --tests  --vendor -e ./...

test: vet ${GOPATH}/bin/gocovmerge
	go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}go test -race -test.v -test.timeout=120s -coverprofile={{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile {{.ImportPath}}{{end}}' $(GOPACKAGES) | xargs -I {} bash -c {}
	@${GOPATH}/bin/gocovmerge `ls *.coverprofile` > cover.out
	@rm -f *.coverprofile
	go tool cover -func=cover.out

coverage: cover.html
cover.html: cover.out
	go tool cover -html=cover.out -o=cover.html
