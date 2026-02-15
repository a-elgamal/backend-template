package telemetry

import "github.com/go-logr/logr"

type logWriter struct {
	logger logr.Logger
	v      int
}

func (w *logWriter) Write(p []byte) (int, error) {
	w.logger.V(w.v).Info(string(p))
	return len(p), nil
}
