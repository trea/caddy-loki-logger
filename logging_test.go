package loki

import (
	"golang.org/x/exp/slices"
	"strings"
	"testing"
)

func TestStringer(t *testing.T) {
	l := LokiLogger{
		endpoint: "https://foo:bar@example.com",
	}

	if strings.Contains(l.String(), "bar") {
		t.Errorf("Expected String() to redact HTTP Basic password: %s", l.String())
	}
}

func TestWriterKeyWithDifferentEndpoints(t *testing.T) {
	tests := []string{
		"http://user:pass@example.com",
		"https://user:pass@example.com",
		"https://example.com",
		"https://127.0.0.1:3000",
		"http://127.0.0.1",
	}

	var prev []string

	logger := LokiLogger{}

	for _, endpoint := range tests {

		logger.endpoint = endpoint
		key := logger.WriterKey()

		if slices.Contains(prev, key) {
			idx := slices.Index(prev, key)
			t.Errorf("Expected WriterKey to be unique, got %s for endpoint %s which is the same as for case: \n %+v", key, endpoint, tests[idx])
			continue
		}

		prev = append(prev, key)
	}
}

func TestWriterKeyWithDifferentLabels(t *testing.T) {
	tests := []map[string]interface{}{
		{
			"foo": "bar",
		},
		{
			"another": "test",
			"var":     "val",
		},
	}

	var prev []string

	logger := LokiLogger{
		endpoint: "https://test.example.com",
	}

	for _, labels := range tests {

		key := logger.WriterKey()

		logger.labels = labels

		if slices.Contains(prev, key) {
			idx := slices.Index(prev, key)
			t.Errorf("Expected WriterKey to be unique, got %s for labels %s which is the same as for case: \n %+v", key, labels, tests[idx])
			continue
		}

		prev = append(prev, key)
	}
}
