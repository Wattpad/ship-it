#!/bin/bash

LINTS=$(go list ./... | grep -v /vendor | xargs -L 1 golint)

if [[ $LINTS ]]; then
    echo "$LINTS"
    exit 1
fi
