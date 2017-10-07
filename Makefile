PKG = github.com/ReformedDevs/anonbot
CMD = anonbot

CWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)

# Find all Go source files (excluding the cache path)
SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')

# Find all resources (static files and templates)
RESOURCES = $(shell find server/static server/templates)

all: dist/${CMD}

# Build the standalone executable
dist/${CMD}: ${SOURCES} server/ab0x.go | cache/lib cache/src/${PKG} dist
	@docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e GIT_COMMITTER_NAME=a \
	    -e GIT_COMMITTER_EMAIL=b \
	    -u ${UID}:${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    golang:latest \
	    go get -pkgdir /go/lib ${PKG}/cmd/${CMD}

# Create a Go source file with the static files and templates
server/ab0x.go: ${RESOURCES} dist/fileb0x b0x.yaml
	@dist/fileb0x b0x.yaml

# Create the fileb0x executable needed for embedding files
dist/fileb0x: | cache/lib cache/src/${PKG} dist
	@docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e GIT_COMMITTER_NAME=a \
	    -e GIT_COMMITTER_EMAIL=b \
	    -u ${UID}:${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    golang:latest \
	    go get -pkgdir /go/lib github.com/UnnoTed/fileb0x

cache/lib:
	@mkdir -p cache/lib

cache/src/${PKG}:
	@mkdir -p cache/src/${PKG}

dist:
	@mkdir dist

clean:
	@rm -rf cache dist server/ab0x.go

.PHONY: clean
