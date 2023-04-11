package loki

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
	"io"
	"net/url"
)

func init() {
	caddy.RegisterModule(LokiLogger{})
}

type LokiLogger struct {
	endpoint string
	labels   map[string]interface{}
	repl     *caddy.Replacer
	logger   *zap.Logger
}

func (l LokiLogger) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.logging.writers.loki",
		New: func() caddy.Module {
			return LokiLogger{}
		},
	}
}

var (
	ErrEmptyEndpoint               = errors.New("endpoint url is required")
	ErrInvalidEndpointPlaceholders = errors.New("invalid endpoint placeholders")
	ErrInvalidEndpointUrl          = errors.New("invalid endpoint url")
	ErrInvalidLabelValReplacement  = errors.New("invalid label value placeholder")
)

func (l *LokiLogger) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if l.labels == nil {
		l.labels = make(map[string]interface{})
	}

	for d.Next() {

		if !d.NextArg() {
			return d.ArgErr()
		}

		l.endpoint = d.Val()

		for nesting := d.Nesting(); d.NextBlock(nesting); {

			if d.Val() == "label" {
				if d.CountRemainingArgs() == 2 {
					args := d.RemainingArgs()
					label := args[0]
					value := args[1]

					l.labels[label] = value

					continue
				}

				for nesting := d.Nesting(); d.NextBlock(nesting); {
					label := d.Val()

					if !d.NextArg() {
						return d.ArgErr()
					}

					l.labels[label] = d.ScalarVal()
				}
			}
		}

	}
	return nil
}

func (l LokiLogger) String() string {
	u, _ := url.Parse(l.endpoint)

	return u.Redacted()
}

func (l LokiLogger) WriterKey() string {
	h := md5.New()
	io.WriteString(h, l.endpoint)
	io.WriteString(h, fmt.Sprintf("%s", l.labels))

	return fmt.Sprintf("loki:%x", h.Sum(nil))
}

func (l LokiLogger) OpenWriter() (io.WriteCloser, error) {
	return NewLokiWriter(l.endpoint, l.labels, l.logger), nil
}

func (l LokiLogger) Validate() error {
	if l.endpoint == "" {
		return ErrEmptyEndpoint
	}

	replacedUrl, err := l.repl.ReplaceOrErr(l.endpoint, true, true)

	if err != nil {
		return errors.Join(ErrInvalidEndpointPlaceholders, err)
	}

	if _, err := url.Parse(replacedUrl); err != nil {
		return errors.Join(ErrInvalidEndpointUrl, err)
	}

	for _, v := range l.labels {
		vStr, ok := v.(string)

		if !ok {
			continue
		}

		_, err := l.repl.ReplaceOrErr(vStr, false, true)

		if err != nil {
			return errors.Join(ErrInvalidLabelValReplacement, err)
		}
	}

	return nil
}

func (l *LokiLogger) Provision(context caddy.Context) error {
	l.repl = context.Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	l.logger = context.Logger(l)
	return nil
}

// Interface guards.
var (
	_ caddy.Module          = (*LokiLogger)(nil)
	_ caddy.Provisioner     = (*LokiLogger)(nil)
	_ caddy.Validator       = (*LokiLogger)(nil)
	_ caddy.WriterOpener    = (*LokiLogger)(nil)
	_ caddyfile.Unmarshaler = (*LokiLogger)(nil)
)
