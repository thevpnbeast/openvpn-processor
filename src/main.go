package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thevpnbeast/openvpn-processor/src/internal/processor"
)

func main() {
	lambda.Start(processor.ProcessEventHandler)
}
