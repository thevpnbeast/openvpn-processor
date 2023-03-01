# Openvpn Processor
[![CI](https://github.com/thevpnbeast/openvpn-processor/workflows/CI/badge.svg?event=push)](https://github.com/thevpnbeast/openvpn-processor/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/thevpnbeast/openvpn-processor)](https://goreportcard.com/report/github.com/thevpnbeast/openvpn-processor)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=thevpnbeast_openvpn-processor&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=thevpnbeast_openvpn-processor)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=thevpnbeast_openvpn-processor&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=thevpnbeast_openvpn-processor)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=thevpnbeast_openvpn-processor&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=thevpnbeast_openvpn-processor)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=thevpnbeast_openvpn-processor&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=thevpnbeast_openvpn-processor)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=thevpnbeast_openvpn-processor&metric=coverage)](https://sonarcloud.io/summary/new_code?id=thevpnbeast_openvpn-processor)
[![Go version](https://img.shields.io/github/go-mod/go-version/thevpnbeast/openvpn-processor)](https://github.com/thevpnbeast/openvpn-processor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Development
This project requires below tools while developing:
- [Golang 1.19](https://golang.org/doc/go1.19)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)
- [gocyclo](https://github.com/fzipp/gocyclo) - required by [pre-commit](https://pre-commit.com/)

Simply run below command to prepare your development environment:
```shell
$ python3 -m venv venv
$ source venv/bin/activate
$ pip3 install pre-commit
$ pre-commit install -c build/ci/.pre-commit-config.yaml
```

Sample SAM commands:
```shell
# Validate the SAM Template
$ make sam-validate
# Invoke function
$ make sam-local-invoke
# Test function in the cloud
$ make sam-cloud-invoke
# Deploy
$ make sam-deploy
```

## Debugging Locally
First, you should spin up the database using [this repository](https://github.com/thevpnbeast/compose).

Then you should provide `docker_network=compose_default` using SAM CLI argument.