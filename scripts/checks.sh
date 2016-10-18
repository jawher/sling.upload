#!/usr/bin/env bash

set -ev

! gofmt -s -d . 2>&1 | read
go test -v
go vet
ineffassign .
! golint . | read