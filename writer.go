package loki

import (
	"fmt"
	"github.com/trea/loki-sink-for-zap"
	"go.uber.org/zap"
	"io"
)

func NewLokiWriter(endpoint string, labels map[string]interface{}, logger *zap.Logger) *LokiWriter {
	rs := &loki_sink_for_zap.LokiWriteSyncer{
		Tags: labels,
		Url:  endpoint,
	}

	return &LokiWriter{
		rs,
		logger,
	}
}

type LokiWriter struct {
	rs     *loki_sink_for_zap.LokiWriteSyncer
	logger *zap.Logger
}

func (l LokiWriter) Write(p []byte) (n int, err error) {
	written, err := l.rs.Write(p)

	if err != nil {
		return written, err
	}

	if err := l.rs.Sync(); err != nil {
		l.logger.Warn(fmt.Sprintf("Writing log entry to Loki failed: %+v", err))
	}

	return written, nil
}

func (l LokiWriter) Close() error {
	return l.rs.Close()
}

var _ io.WriteCloser = (*LokiWriter)(nil)
