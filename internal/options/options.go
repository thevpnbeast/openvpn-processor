package options

import (
	commons "github.com/thevpnbeast/golang-commons"
	"go.uber.org/zap"
)

const (
	ApplicationName      = "openvpn-processor"
	TargetConfigFileName = "application"
	TargetConfigFileType = "yaml"
)

var (
	opts   *OpenvpnProcessorOptions
	logger *zap.Logger
)

func init() {
	applicationId := getStringEnv("CONFIG_APPLICATION_ID", "jl9vegs")
	confProfileId := getStringEnv("CONFIG_PROFILE_ID", "hek207s")

	logger = commons.GetLogger()
	opts = newOpenvpnProcessorOptions()

	if err := fetchConfig(applicationId, confProfileId); err != nil {
		logger.Fatal("fatal error occured while fetching config from AWS AppConfig", zap.String("error", err.Error()))
	}
}

// OpenvpnProcessorOptions represents openvpn-processor environment variables
type OpenvpnProcessorOptions struct {
	VpnGateUrl            string
	DbUrl                 string
	DbDriver              string
	DbMaxOpenConn         int
	DbMaxIdleConn         int
	DbConnMaxLifetimeMin  int
	DialTcpTimeoutSeconds int
}

// GetOpenvpnProcessorOptions returns the initialized VpnbeastServiceOptions
func GetOpenvpnProcessorOptions() *OpenvpnProcessorOptions {
	return opts
}

// newAuthServiceOptions creates an AuthServiceOptions struct with zero values
func newOpenvpnProcessorOptions() *OpenvpnProcessorOptions {
	return &OpenvpnProcessorOptions{}
}
