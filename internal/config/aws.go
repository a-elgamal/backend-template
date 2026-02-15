package config

import "github.com/spf13/viper"

const awsALBRegion = "AWS.ALB_REGION"

// AWSConfig contains AWS-specific configuration
type AWSConfig struct {
	v *viper.Viper
}

// ALBAuthEnabled reports whether ALB OIDC authentication is enabled
func (c AWSConfig) ALBAuthEnabled() bool {
	return c.v.GetString(awsALBRegion) != ""
}

// ALBRegion returns the AWS region for ALB public key retrieval
func (c AWSConfig) ALBRegion() string {
	return c.v.GetString(awsALBRegion)
}
