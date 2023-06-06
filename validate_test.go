package loki

import (
	"context"
	"errors"
	"github.com/caddyserver/caddy/v2"
	"os"
	"testing"
)

func TestEmptyEndpointUrl(t *testing.T) {
	l := LokiLogger{}
	replacer := caddy.NewReplacer()
	base := context.WithValue(context.Background(), caddy.ReplacerCtxKey, replacer)
	ctx, cancel := caddy.NewContext(caddy.Context{Context: base})
	defer cancel()

	if err := l.Provision(ctx); err != nil {
		t.Fatalf("Expected provision to pass, got err: %+v", err)
	}

	err := l.Validate()

	if err == nil {
		t.Errorf("Expected validation to fail for empty endpoint URL")
	}

	if !errors.Is(err, ErrEmptyEndpoint) {
		t.Errorf("Expected err to be ErrEmptyEndpoint, got: %+v", err)
	}
}

func TestInvalidEndpointPlaceholders(t *testing.T) {
	l := LokiLogger{}
	l.Endpoint = "{vars.nonExistentUrl}"
	replacer := caddy.NewReplacer()
	base := context.WithValue(context.Background(), caddy.ReplacerCtxKey, replacer)
	ctx, cancel := caddy.NewContext(caddy.Context{Context: base})
	defer cancel()

	if err := l.Provision(ctx); err != nil {
		t.Fatalf("Expected provision to pass, got err: %+v", err)
	}
	err := l.Validate()

	if err == nil {
		t.Errorf("Expected validation to fail for invalid placeholder")
		return
	}

	if !errors.Is(err, ErrInvalidEndpointPlaceholders) {
		t.Errorf("Expected err to be ErrInvalidEndpointPlaceholders, got: %+v", err)
	}
}

func TestGoodEndpointPlaceholder(t *testing.T) {
	os.Setenv("LOKI_ENDPOINT", "https://test:password@example.com")
	t.Cleanup(func() {
		os.Unsetenv("LOKI_ENDPOINT")
	})

	l := LokiLogger{}
	l.Endpoint = "{env.LOKI_ENDPOINT}"

	replacer := caddy.NewReplacer()
	base := context.WithValue(context.Background(), caddy.ReplacerCtxKey, replacer)
	ctx, cancel := caddy.NewContext(caddy.Context{Context: base})
	defer cancel()

	if err := l.Provision(ctx); err != nil {
		t.Fatalf("Expected provision to pass, got err: %+v", err)
	}

	err := l.Validate()

	if err != nil {
		t.Errorf("Expected validation to pass, got err: %+v", err)
	}
}

func TestLabelPlaceholders(t *testing.T) {
	l := LokiLogger{
		Endpoint: "https://example.com",
		Labels: map[string]interface{}{
			"app_env":        "{env.APP_ENV}",
			"varplaceholder": "{somevar}",
		},
	}

	replacer := caddy.NewReplacer()
	replacer.Set("somevar", "val")
	base := context.WithValue(context.Background(), caddy.ReplacerCtxKey, replacer)
	ctx, cancel := caddy.NewContext(caddy.Context{Context: base})
	defer cancel()

	if err := l.Provision(ctx); err != nil {
		t.Fatalf("Expected provision to pass, got err: %+v", err)
	}

	err := l.Validate()

	if err != nil {
		t.Errorf("expected validation to pass, got err: %+v", err)
	}
}
