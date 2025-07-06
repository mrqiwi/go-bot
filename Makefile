PROJECTNAME=$(shell basename "$(PWD)")
BUILD_DIR=build/package/

.PHONY: build
build:
	go build -o ${BUILD_DIR}${PROJECTNAME} cmd/app/main.go

.PHONY: build-arm
build-arm:
	GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}${PROJECTNAME}-arm cmd/app/main.go

.PHONY: run
run: build
	./${BUILD_DIR}${PROJECTNAME}

.PHONY: update
update: build
	cp ./${BUILD_DIR}${PROJECTNAME} /usr/sbin/gobot