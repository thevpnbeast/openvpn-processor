package options

import (
	"log"
)

const (
	ApplicationName      = "openvpn-processor"
	TargetConfigFileName = "application"
	TargetConfigFileType = "yaml"
)

var opts *OpenvpnProcessorOptions

func init() {
	applicationId := getStringEnv("CONFIG_APPLICATION_ID", "jl9vegs")
	confProfileId := getStringEnv("CONFIG_PROFILE_ID", "hek207s")

	opts = newOpenvpnProcessorOptions()

	if err := fetchConfig(applicationId, confProfileId); err != nil {
		log.Printf("FATAL: fatal error occurred while fetching config from AWS AppConfig (error=%s)", err.Error())
		panic(err)
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
