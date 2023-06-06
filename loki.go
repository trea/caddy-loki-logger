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
	Endpoint string
	Labels   map[string]interface{}
	Repl     *caddy.Replacer
	Logger   *zap.Logger
}

func (l LokiLogger) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.logging.writers.loki",
		New: func() caddy.Module {
			return new(LokiLogger)
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
	if l.Labels == nil {
		l.Labels = make(map[string]interface{})
	}

	for d.Next() {

		if !d.NextArg() {
			return d.ArgErr()
		}

		l.Endpoint = d.Val()

		for nesting := d.Nesting(); d.NextBlock(nesting); {

			if d.Val() == "label" {
				if d.CountRemainingArgs() == 2 {
					args := d.RemainingArgs()
					label := args[0]
					value := args[1]

					l.Labels[label] = value

					continue
				}

				for nesting := d.Nesting(); d.NextBlock(nesting); {
					label := d.Val()

					if !d.NextArg() {
						return d.ArgErr()
					}

					l.Labels[label] = d.ScalarVal()
				}
			}
		}

	}
	return nil
}

func (l LokiLogger) String() string {
	u, _ := url.Parse(l.Endpoint)

	return u.Redacted()
}

func (l LokiLogger) WriterKey() string {
	h := md5.New()
	io.WriteString(h, l.Endpoint)
	io.WriteString(h, fmt.Sprintf("%s", l.Labels))

	return fmt.Sprintf("loki:%x", h.Sum(nil))
}

func (l *LokiLogger) OpenWriter() (io.WriteCloser, error) {
	return NewLokiWriter(l.Endpoint, l.Labels, l.Logger), nil
}

func (l *LokiLogger) Validate() error {
	if l.Endpoint == "" {
		return ErrEmptyEndpoint
	}

	replacedUrl, err := l.Repl.ReplaceOrErr(l.Endpoint, true, true)

	if err != nil {
		return errors.Join(ErrInvalidEndpointPlaceholders, err)
	}

	if _, err := url.Parse(replacedUrl); err != nil {
		return errors.Join(ErrInvalidEndpointUrl, err)
	}

	for _, v := range l.Labels {
		vStr, ok := v.(string)

		if !ok {
			continue
		}

		_, err := l.Repl.ReplaceOrErr(vStr, false, true)

		if err != nil {
			return errors.Join(ErrInvalidLabelValReplacement, err)
		}
	}

	return nil
}

func (l *LokiLogger) Provision(context caddy.Context) error {
	repl, ok := context.Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	if !ok {
		//return fmt.Errorf("unable to get caddy replacer")
		repl = caddy.NewReplacer()
	}

	l.Repl = repl
	l.Logger = context.Logger(l)
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
