#!/bin/bash

# Rerun the tests whenever a file changes in this project.
#
# You must install `ag` ("The Silver Searcher") and `entr` to use this
# script.
#
# See: http://eradman.com/entrproject/
# See: https://github.com/ggreer/the_silver_searcher

while true; do
    ag -l --go | entr -c -d go test ./... && break
done
