#!/bin/bash

## Script to release grepj
## Build Go project for multiple architectures - http://golang.org/doc/install/source
## Publish to s3


OUT=grepj
SRC="$OUT.go"
DIST=release
S3_BUCKET=s3://grepj

declare -r PLATFORMS=(
    "darwin"
    "darwin"
    "linux"
    "linux"
    "windows"
    "windows"
)

declare -r ARCHS=(
    "amd64"
    "386"
    "amd64"
    "386"
    "amd64"
    "386"
)

declare -r EXTENSIONS=(
    ""
    ""
    ""
    ""
    ".exe"
    ".exe"
)

function clean {
    rm -rf $DIST
    echo "Clean"
}

function build_one {
    GOOS=$1 GOARCH=$2 go build -o "$DIST/$1_$2/$OUT$3" $SRC
    echo "Built for $1 $2"
}

function build {
    i=0
    while [ $i -lt ${#PLATFORMS[*]} ]; do
	build_one ${PLATFORMS[$i]} ${ARCHS[$i]} ${EXTENSIONS[$i]}
	i=$(( $i + 1 ))
    done
}

function release_one {
    s3cmd put --acl-public --add-header "Cache-Control: no-cache" \
	"$DIST/$1_$2/$OUT$3" "$S3_BUCKET/$DIST/$1_$2/$OUT$3"
    echo "Release $1 $2 $3"
}

function release {
    i=0
    while [ $i -lt ${#PLATFORMS[*]} ]; do
        release_one ${PLATFORMS[$i]} ${ARCHS[$i]} ${EXTENSIONS[$i]}
        i=$(( $i + 1 ))
    done
}

clean
build
release
