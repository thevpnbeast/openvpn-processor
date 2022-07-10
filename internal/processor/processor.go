package processor

import (
	commons "github.com/thevpnbeast/golang-commons"
	"github.com/thevpnbeast/openvpn-processor/internal/options"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	opts   *options.OpenvpnProcessorOptions
)

func init() {
	logger = commons.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
}

func ProcessEventHandler() {
	// TODO: implement
	logger.Info("starting processEventHandler")
}
