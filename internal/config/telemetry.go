package config

import (
	"strings"

	"github.com/spf13/viper"
)

const tracingConfigSampling = "TELEMETRY.TRACING.SAMPLING"
const loggingConfigLevel = "TELEMETRY.LOGGING.LEVEL"
const loggingSyslogTCPAddress = "TELEMETRY.LOGGING.SYSLOG_TCP_ADDRESS"
const loggingConsoleLoggingEnabled = "TELEMETRY.LOGGING.CONSOLE_LOGGING_ENABLED"

// TelemetryConfig contains all the telemetry-related configuration of the application
type TelemetryConfig struct {
	v *viper.Viper
}

// NewTelemetryConfig creates a new Tracing Config directly
func NewTelemetryConfig(loggingLevel string, syslogTCPAddress string, traceSampling float64) TelemetryConfig {
	v := viper.New()
	v.Set(loggingConfigLevel, loggingLevel)
	v.Set(tracingConfigSampling, traceSampling)
	v.Set(loggingSyslogTCPAddress, syslogTCPAddress)
	return TelemetryConfig{v: v}
}

// TraceSampling returns a number between 0->1 that controls the tracing sampling ratio
func (tc TelemetryConfig) TraceSampling() float64 {
	return tc.v.GetFloat64(tracingConfigSampling)
}

// LogLevel indicates the minimum zerlog level at which the logging should happen
func (tc TelemetryConfig) LogLevel() string {
	return tc.v.GetString(loggingConfigLevel)
}

// SyslogTCPEnabled determines whether Syslog over TCP should be turned or not
func (tc TelemetryConfig) SyslogTCPEnabled() bool {
	return tc.v.IsSet(loggingSyslogTCPAddress) && strings.TrimSpace(tc.v.GetString(loggingSyslogTCPAddress)) != ""
}

// SyslogTCPAddress the network address to send syslog data to over TCP
func (tc TelemetryConfig) SyslogTCPAddress() string {
	return tc.v.GetString(loggingSyslogTCPAddress)
}

// ConsoleWriterEnabled Whether logging to console should be turned on or not
func (tc TelemetryConfig) ConsoleWriterEnabled() bool {
	return tc.v.GetBool(loggingConsoleLoggingEnabled)
}
