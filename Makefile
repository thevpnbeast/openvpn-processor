# assumes that we have already a profile named thevpnbeast-root in AWS CLI config
#export AWS_PROFILE := thevpnbeast-root
#
APP_NAME = openvpn-processor
AWS_REGION = us-east-1
AWS_IAM_CAPABILITIES = CAPABILITY_IAM
AWS_RELEASES_BUCKET = thevpnbeast-releases-1
AWS_STACK_NAME = openvpn-processor
TEMPLATE_FILE = template.yaml
GENERATED_TEMPLATE_FILE = template_generated.yaml

ERRCHECK_VERSION = latest
GOLANGCI_LINT_VERSION = latest
REVIVE_VERSION = latest
GOIMPORTS_VERSION = latest
INEFFASSIGN_VERSION = latest

LOCAL_BIN := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/.bin
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "1.0.0")

.PHONY: all
all: clean tools lint fmt test build

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)

.PHONY: tools
tools:  golangci-lint-install revive-install go-imports-install ineffassign-install
	go mod tidy

.PHONY: golangci-lint-install
golangci-lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: revive-install
revive-install:
	GOBIN=$(LOCAL_BIN) go install github.com/mgechev/revive@$(REVIVE_VERSION)

.PHONY: ineffassign-install
ineffassign-install:
	GOBIN=$(LOCAL_BIN) go install github.com/gordonklaus/ineffassign@$(INEFFASSIGN_VERSION)

.PHONY: lint
lint: tools run-lint

.PHONY: run-lint
run-lint: lint-golangci-lint lint-revive

.PHONY: lint-golangci-lint
lint-golangci-lint:
	$(info running golangci-lint...)
	$(LOCAL_BIN)/golangci-lint -v run ./... || (echo golangci-lint returned an error, exiting!; sh -c 'exit 1';)

.PHONY: lint-revive
lint-revive:
	$(info running revive...)
	$(LOCAL_BIN)/revive -formatter=stylish -config=build/ci/.revive.toml -exclude ./vendor/... ./... || (echo revive returned an error, exiting!; sh -c 'exit 1';)

.PHONY: upgrade-direct-deps
upgrade-direct-deps: tidy
	for item in `grep -v 'indirect' go.mod | grep '/' | cut -d ' ' -f 1`; do \
		echo "trying to upgrade direct dependency $$item" ; \
		go get -u $$item ; \
  	done
	go mod tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

.PHONY: run-goimports
run-goimports: go-imports-install
	for item in `find . -type f -name '*.go' -not -path './vendor/*'`; do \
		$(LOCAL_BIN)/goimports -l -w $$item ; \
	done

.PHONY: go-imports-install
go-imports-install:
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

.PHONY: fmt
fmt: tools run-fmt run-ineffassign run-vet

.PHONY: run-fmt
run-fmt:
	$(info running fmt...)
	go fmt ./... || (echo fmt returned an error, exiting!; sh -c 'exit 1';)

.PHONY: run-ineffassign
run-ineffassign:
	$(info running ineffassign...)
	$(LOCAL_BIN)/ineffassign ./... || (echo ineffassign returned an error, exiting!; sh -c 'exit 1';)

.PHONY: run-vet
run-vet:
	$(info running vet...)
	go vet ./... || (echo vet returned an error, exiting!; sh -c 'exit 1';)

.PHONY: test
test: tidy
	$(info starting the test for whole module...)
	go test -failfast -vet=off -race ./... || (echo an error while testing, exiting!; sh -c 'exit 1';)

.PHONY: test-with-coverage
test-with-coverage: tidy
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: update
update: tidy
	go get -u ./...

.PHONY: build
build: tidy
	$(info building binary...)
	GOOS=linux GOARCH=amd64 go build -o src/main src/main.go || (echo an error while building binary, exiting!; sh -c 'exit 1';)

.PHONY: run
run: tidy
	go run src/main.go

.PHONY: cross-compile
cross-compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 src/main.go
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 src/main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 src/main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 src/main.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 src/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 src/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 src/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 src/main.go

.PHONY: aws-build
aws-build:
	go get -v all
	GOOS=linux GOARCH=amd64 go build -o src/main src/main.go
	zip -jrm src/main-$(VERSION).zip src/main


.PHONY: aws-deploy
aws-deploy: aws-build
	aws lambda update-function-code --function-name openvpn-processor --zip-file fileb://src/main.zip

.PHONY: aws-publish
aws-publish: aws-build
	aws lambda update-function-code --function-name openvpn-processor --zip-file fileb://src/main.zip --publish

.PHONY: sam-validate
sam-validate:
	sam validate

.PHONY: sam-local-start-api
sam-local-start-api:
	sam local start-api

.PHONY: sam-local-invoke
sam-local-invoke:
	sam local invoke

.PHONY: sam-cloud-invoke
sam-cloud-invoke:
	sam sync --stack-name $(AWS_STACK_NAME) --watch

.PHONY: sam-build
sam-build: build
	which build-lambda-zip || go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest
	build-lambda-zip -o src/main.zip src/main || (echo an error while compressing binary with build-lambda-zip, exiting!; sh -c 'exit 1';)
	sam build

.PHONY: sam-package
sam-package: sam-build
	sam package --s3-prefix $(APP_NAME) --s3-bucket $(AWS_RELEASES_BUCKET) --template-file $(TEMPLATE_FILE) --output-template-file $(GENERATED_TEMPLATE_FILE)

.PHONY: sam-deploy
sam-deploy: sam-package
	sam deploy --parameter-overrides StageName=prod AppVersion=$(VERSION) --no-confirm-changeset \
		--no-fail-on-empty-changeset --template-file $(GENERATED_TEMPLATE_FILE) --stack-name $(AWS_STACK_NAME) \
		--s3-prefix $(APP_NAME) --s3-bucket $(AWS_RELEASES_BUCKET) --capabilities $(AWS_IAM_CAPABILITIES) \
		--region $(AWS_REGION)

.PHONY: sam-publish
sam-publish: sam-deploy
	sam publish --region $(AWS_REGION) --semantic-version $(VERSION)
