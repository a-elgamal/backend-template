package telemetry

import (
	"io"
	"log/syslog"
	"os"
	"strings"
	"time"

	"alielgamal.com/myservice/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

const syslogTag = "myservice"
const syslogPriority = syslog.LOG_DEBUG

// SetupLogger Returns a configured logr.Logger instance
func SetupLogger(config config.TelemetryConfig) (logr.Logger, error) {
	zerolog.TimeFieldFormat = time.StampNano
	lvl, err := zerolog.ParseLevel(strings.ToLower(config.LogLevel()))
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	// Logging will always be set as one of Info, Debug, or Trace. Other levels map to Info
	zerologr.SetMaxV(1 - int(lvl))

	writers := []io.Writer{}
	if config.ConsoleWriterEnabled() {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampNano})
	}
	if sl, err := syslog.New(syslogPriority, syslogTag); err != nil {
		log.Err(err).Msg("Local Syslog writer failed")
	} else {
		writers = append(writers, zerolog.SyslogCEEWriter(sl))
	}

	if config.SyslogTCPEnabled() {
		if sl, err := syslog.Dial("tcp", config.SyslogTCPAddress(), syslogPriority, syslogTag); err != nil {
			log.Err(err).Msg("TCP Syslog writer failed")
		} else {
			writers = append(writers, zerolog.SyslogCEEWriter(sl))
		}
	}

	multiWriter := zerolog.MultiLevelWriter(writers...)

	zl := zerolog.New(multiWriter).Level(zerolog.GlobalLevel()).With().Timestamp().Logger()
	var logger logr.Logger = zerologr.New(&zl)
	gin.DisableConsoleColor()
	gin.DefaultWriter = &logWriter{logger, 1}
	gin.DefaultErrorWriter = &logWriter{logger, 0}
	otel.SetLogger(logger.WithName("OpenTelemetry"))

	return logger, nil
}
