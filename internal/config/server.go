package config

import "github.com/spf13/viper"

const serverConfigHTTPAddress = "SERVER.HTTP_ADDRESS"
const serverCORSAllowedOrigins = "SERVER.CORS_ALLOWED_ORIGINS"
const serverPortalPath = "SERVER.PORTAL_PATH"

// ServerConfig contains all the configurations related to the server
type ServerConfig struct {
	v *viper.Viper
}

// GetHTTPAddress contains the address that the server should listen to
func (sc ServerConfig) GetHTTPAddress() string {
	return sc.v.GetString(serverConfigHTTPAddress)
}

// CORSAllowedOrigins The list of origins that should be allowed in CORS policy.
func (sc ServerConfig) CORSAllowedOrigins() []string {
	return sc.v.GetStringSlice(serverCORSAllowedOrigins)
}

// PortalPath The path from which portal files should be served
func (sc ServerConfig) PortalPath() string {
	return sc.v.GetString(serverPortalPath)
}
