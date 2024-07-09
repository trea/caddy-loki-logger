package loki

import (
	"context"
	"fmt"
	"github.com/trea/loki-sink-for-zap"
	"go.uber.org/zap"
	"io"
	"log"
)

func NewLokiWriter(endpoint string, labels map[string]interface{}, logger *zap.Logger) *LokiWriter {
	ctx, cancel := context.WithCancel(context.Background())

	ws := loki_sink_for_zap.NewLokiWriteSyncer(ctx)
	ws.Url = endpoint
	ws.Tags = labels

	defer cancel()

	return &LokiWriter{
		rs:     ws,
		logger: logger,
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
		log.Printf(fmt.Sprintf("Writing log entry to Loki failed: %+v", err))
	}

	return written, nil
}

func (l LokiWriter) Close() error {
	return l.rs.Close()
}

var _ io.WriteCloser = (*LokiWriter)(nil)
