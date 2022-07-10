package processor

import (
	commons "github.com/thevpnbeast/golang-commons"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = commons.GetLogger()
}

func ProcessEventHandler() error {
	// TODO: implement
	logger.Info("starting processEventHandler")
	return nil
}
