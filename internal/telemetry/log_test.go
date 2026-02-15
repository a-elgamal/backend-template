package telemetry

import (
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"alielgamal.com/myservice/internal/config"
)

func TestSetupLogging(t *testing.T) {

	t.Run("Properly configures the logger's level based on the config", func(t *testing.T) {

		assertLevel := func(t *testing.T, levelConfig string, expectedLevel zerolog.Level) {
			tc := config.NewTelemetryConfig(levelConfig, "", 0)

			_, err := SetupLogger(tc)
			assert.NoError(t, err)
			assert.Equal(t, expectedLevel, zerolog.GlobalLevel())
		}

		cases := []struct {
			name   string
			config string
			level  zerolog.Level
		}{
			{"debug", "debug", zerolog.DebugLevel},
			{"info", "info", zerolog.InfoLevel},
			{"warn", "warn", zerolog.InfoLevel},
			{"error", "error", zerolog.InfoLevel},
			{"fatal", "fatal", zerolog.InfoLevel},
			{"panic", "panic", zerolog.InfoLevel},
			{"No Level", "", zerolog.InfoLevel},
			{"disabled", "disabled", zerolog.InfoLevel},
			{"trace", "trace", zerolog.TraceLevel},
			{"Case Insensitive Config", "tRaCe", zerolog.TraceLevel},
			{"Invalid Level", "Invalid", zerolog.InfoLevel},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				assertLevel(t, tt.config, tt.level)
			})
		}
	})

	t.Run("Sets up Syslog TCP Logging when enabled", func(t *testing.T) {

		listener, err := net.Listen("tcp", ":0")
		assert.NoError(t, err)
		defer listener.Close()

		connected := false
		go func() {
			conn, _ := listener.Accept()
			defer conn.Close()
			connected = true
		}()

		tc := config.NewTelemetryConfig("debug", listener.Addr().String(), 0)

		_, err = SetupLogger(tc)
		assert.NoError(t, err)

		assert.Eventually(t, func() bool { return connected }, time.Minute, time.Second, "TCP Syslog connection not established")
	})
}
