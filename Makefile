PROJECTNAME=$(shell basename "$(PWD)")
BUILD_DIR=build/package/

.PHONY: build
build:
	go build -o ${BUILD_DIR}${PROJECTNAME} cmd/app/main.go

.PHONY: run
run:
	./${BUILD_DIR}${PROJECTNAME}