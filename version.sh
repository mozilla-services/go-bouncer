#!/bin/bash

set -e

cd "$(dirname "$0")"

VERSION=1.0.0
if REV=$(git rev-parse --short HEAD); then
    VERSION="${VERSION}-${REV}"
fi


cat > bouncer/version.go <<HERE
package bouncer

const Version = "${VERSION}"
HERE
