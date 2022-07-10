package options

var opts *VpnbeastServiceOptions

// VpnbeastServiceOptions represents openvpn-processor environment variables
type VpnbeastServiceOptions struct {
	DbUrl                string `env:"DB_URL"`
	DbDriver             string `env:"DB_DRIVER"`
	DbMaxOpenConn        int    `env:"DB_MAX_OPEN_CONN"`
	DbMaxIdleConn        int    `env:"DB_MAX_IDLE_CONN"`
	DbConnMaxLifetimeMin int    `env:"DB_CONN_MAX_LIFETIME_MIN"`
}

// GetVpnbeastServiceOptions returns the initialized VpnbeastServiceOptions
func GetVpnbeastServiceOptions() *VpnbeastServiceOptions {
	return opts
}

// newAuthServiceOptions creates an AuthServiceOptions struct with zero values
func newVpnbeastServiceOptions() *VpnbeastServiceOptions {
	return &VpnbeastServiceOptions{}
}
