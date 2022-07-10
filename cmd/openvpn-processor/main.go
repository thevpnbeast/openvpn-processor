package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	commons "github.com/thevpnbeast/golang-commons"
	"github.com/thevpnbeast/openvpn-processor/internal/processor"
)

var logger = commons.GetLogger()

func main() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	lambda.Start(processor.ProcessEventHandler)
}
