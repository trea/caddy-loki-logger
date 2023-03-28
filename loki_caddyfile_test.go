package loki

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"testing"
)

func testEndpoint(t *testing.T, endpoint, expected string) {
	if endpoint != expected {
		t.Errorf("Expected endpoint to be set as '%s' got '%s", endpoint, expected)
	}
}

func testLabel(t *testing.T, labels map[string]interface{}, name string, expected interface{}) {
	val, ok := labels[name]

	if !ok {
		t.Errorf("Expected label %s to have value %+v, instead it wasn't set at all", name, expected)
		return
	}

	if val != expected {
		t.Errorf("Expected label %s to have value %T %+v, instead it had value %T %+v", name, expected, expected, val, val)
	}
}

func testNoLabelsSet(t *testing.T, labels map[string]interface{}) {
	if len(labels) != 0 {
		t.Errorf("Expected no labels to be set, found: %+v", labels)
	}
}

func TestUnmarshalInvalidCaddyfile(t *testing.T) {
	l := &LokiLogger{}
	cf := `
loki
`
	d := caddyfile.NewTestDispenser(cf)
	err := l.UnmarshalCaddyfile(d)

	if err == nil {
		t.Error("Unmarshal passed, should have failed because of no endpoint set")
	}

	testNoLabelsSet(t, l.labels)
}

func TestUnmarshalBasicCaddyfile(t *testing.T) {
	endpoint := "https://myuser:somepassword@example.net/loki/api/v1/push"

	l := &LokiLogger{}
	cf := "loki " + endpoint
	d := caddyfile.NewTestDispenser(cf)
	err := l.UnmarshalCaddyfile(d)

	if err != nil {
		t.Errorf("Unmarshal failed, should have passed: %+v", err)
	}

	testEndpoint(t, l.endpoint, endpoint)
	testNoLabelsSet(t, l.labels)
}

func TestUnmarshalReplacedEnpoint(t *testing.T) {
	endpoint := "{env.LOKI_ENDPOINT}"
	l := &LokiLogger{}
	cf := "loki " + endpoint
	d := caddyfile.NewTestDispenser(cf)
	err := l.UnmarshalCaddyfile(d)

	if err != nil {
		t.Errorf("Unmarshal failed, should have passed: %+v", err)
	}

	testEndpoint(t, l.endpoint, endpoint)
	testNoLabelsSet(t, l.labels)
}

func TestUnmarshalCaddyfileWithStaticLabel(t *testing.T) {
	endpoint := "https://myuser:somepassword@example.net/loki/api/v1/push"

	l := &LokiLogger{}
	cf := "loki " + endpoint
	d := caddyfile.NewTestDispenser(cf)
	err := l.UnmarshalCaddyfile(d)

	if err != nil {
		t.Errorf("Unmarshal failed, should have passed: %+v", err)
	}

	testEndpoint(t, l.endpoint, endpoint)
}

func TestUnmarshalCaddyfileWithNestedLabels(t *testing.T) {
	endpoint := "https://myuser:somepassword@example.net/loki/api/v1/push"
	l := &LokiLogger{}
	cf := "loki " + endpoint + " {\n" +
		`label {
		somelabel example
		anotherlabel 1
		testing	true
		anEnvVar {env.SOMEENV}
		aVar {vars.somevar}
	}
}`

	d := caddyfile.NewTestDispenser(cf)
	err := l.UnmarshalCaddyfile(d)

	if err != nil {
		t.Errorf("Unmarshal failed, should have passed: %+v", err)
	}

	testEndpoint(t, l.endpoint, endpoint)

	testLabel(t, l.labels, "somelabel", "example")
	testLabel(t, l.labels, "anotherlabel", 1)
	testLabel(t, l.labels, "testing", true)
	testLabel(t, l.labels, "anEnvVar", "{env.SOMEENV}")
	testLabel(t, l.labels, "aVar", "{vars.somevar}")
}
