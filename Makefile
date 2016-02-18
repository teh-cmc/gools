SHELL=/bin/bash

NAME=gools
SUBPACKAGES=$(shell go list ./... | grep -v /vendor/)

help:
	@echo "Available targets"
	@echo "================="
	@echo "vet:   runs vetting and linting tools on all of ${NAME}'s subpackages"

vet:
	@GO15VENDOREXPERIMENT=1 gometalinter -E vet -E vetshadow -E golint -E errcheck -E gotype -E structcheck -E deadcode -E dupl -E interfacer -E ineffassign -E varcheck -E aligncheck ./...
