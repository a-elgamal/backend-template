package config

import "github.com/spf13/viper"

const gcpProjectNumber = "GCP.PROJECT_NUMBER"
const gcpRegion = "GCP.REGION"
const gcpInternalBackendServiceID = "GCP.INTERNAL_BACKEND_SERVICE_ID"

// GCPConfig contains GCP-specific configuration
type GCPConfig struct {
	v *viper.Viper
}

// ProjectNumber The GCP project number
func (c GCPConfig) ProjectNumber() int64 {
	return c.v.GetInt64(gcpProjectNumber)
}

// IAPAuthEnabled Whether IAP Authentication is enabled for internal backend service
func (c GCPConfig) IAPAuthEnabled() bool {
	return c.v.IsSet(gcpInternalBackendServiceID)
}

// InternalBackendServiceID The id of the backend service that is configured to handle internal
func (c GCPConfig) InternalBackendServiceID() int64 {
	return c.v.GetInt64(gcpInternalBackendServiceID)
}

// Region The region in which the backend service is deployed
func (c GCPConfig) Region() string {
	return c.v.GetString(gcpRegion)
}
