# assumes that we have already a profile named thevpnbeast-root in AWS CLI config
export AWS_PROFILE := thevpnbeast-root
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0
export GO111MODULE := on

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

ineffassign:
	go get github.com/gordonklaus/ineffassign
	go mod vendor
	ineffassign ./...

test:
	go test ./...

test_coverage:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

build:
	go build -o bin/main cmd/openvpn-processor/main.go

run:
	go run cmd/openvpn-processor/main.go

cross-compile:
	# 32-Bit Systems
	# FreeBDS
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 cmd/openvpn-processor/main.go
	# MacOS
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 cmd/openvpn-processor/main.go
	# Linux
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 cmd/openvpn-processor/main.go
	# Windows
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 cmd/openvpn-processor/main.go
        # 64-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 cmd/openvpn-processor/main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 cmd/openvpn-processor/main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 cmd/openvpn-processor/main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 cmd/openvpn-processor/main.go

upgrade-direct-deps:
	for item in `grep -v 'indirect' go.mod | grep '/' | cut -d ' ' -f 1`; do \
		echo "trying to upgrade direct dependency $$item" ; \
		go get -u $$item ; \
  	done
	go mod tidy
	go mod vendor

aws_build:
	go get -v all
	GOOS=linux go build -o bin/main cmd/openvpn-processor/main.go
	zip -jrm bin/main.zip bin/main

aws_upload: aws_build
	aws lambda update-function-code --function-name openvpn-processor --zip-file fileb://bin/main.zip

aws_upload_publish: aws_build
	aws lambda update-function-code --function-name openvpn-processor --zip-file fileb://bin/main.zip --publish

all: test build run
