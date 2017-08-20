#!/bin/bash

# go get golang.org/x/tools/cmd/goimports

gofmt -s -w app.go
goimports -w app.go