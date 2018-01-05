#!/usr/bin/env bash
set -e
echo "" > coverage.txt
# waiting for oak v1.5.0
# go test -coverprofile=profile.out -covermode=atomic ./unique
# if [ -f profile.out ]; then
#     cat profile.out >> coverage.txt
#     rm profile.out
# fi
go test -coverprofile=profile.out -covermode=atomic ./selection
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./pop
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./pairing
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./neural
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./mut
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./mut/frange
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./mut/irange
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./mut/mutenv
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./lgp
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./gp
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
# waiting on oak v2.0.0
# go test -coverprofile=profile.out -covermode=atomic ./gg
# if [ -f profile.out ]; then
#     cat profile.out >> coverage.txt
#     rm profile.out
# fi
# go test -coverprofile=profile.out -covermode=atomic ./gg/dev
# if [ -f profile.out ]; then
#     cat profile.out >> coverage.txt
#     rm profile.out
# fi
# go test -coverprofile=profile.out -covermode=atomic ./gg/player
# if [ -f profile.out ]; then
#     cat profile.out >> coverage.txt
#     rm profile.out
# fi
go test -coverprofile=profile.out -covermode=atomic ./gevaerr
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./env
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./eda
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./eda/fitness
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./eda/stat
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./cross
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi