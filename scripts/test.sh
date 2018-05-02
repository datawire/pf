#!/bin/bash
set -e

clean() {
    sudo pfctl -F all || true
    sudo pfctl -a asd -F all || true
    sudo pfctl -a myanchor -F all || true
    sudo pfctl -d || true
}


clean

{
    go test -exec sudo -cover -v github.com/datawire/pf
    RESULT=$?
} || true

clean

exit $RESULT
