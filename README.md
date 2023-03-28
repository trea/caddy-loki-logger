# caddy-loki-logger

Push logs from your [Caddy](https://caddyserver.com) web server directly to a [Loki](http://github.com/grafanalabs/loki) instance.

## Usage

### Building

To use this module with Caddy, you'll need to create a custom build with [`xcaddy`](https://github.com/caddyserver/xcaddy).

```bash
xcaddy build --with github.com/trea/caddy-loki-logger
```

#### Via Docker

You can build this module (and others) into Caddy with `xcaddy` using Docker. The following example uses a multi-phase build to build the
Caddy binary with xcaddy, and then copies your newly built binary over the original binary in the core Caddy image.

```Dockerfile
FROM caddy:2.6.4-builder AS builder

RUN xcaddy --with github.com/trea/caddy-loki-logger

FROM caddy:2.6.4

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
```

<!-- TODO: Document Caddy download page https://caddyserver.com/download -->

## Caddyfile Configuration Examples

## Accepted Caddyfile Syntax

```
(
    log default {
        output loki <loki_endpoint> {
            label {
                <label_name> <label_value>
                ...
            }
        }
    }
)
```

`loki_endpoint` can be any valid HTTP/HTTPS URI, including using basic authentication. It may also contain variables including environment variables like `{env.LOKI_ENDPOINT}`.

`label` is not required. It may contain as many label-value pairs as desired. `label_value` may contain variables including environment variables.

The `output` does not have to be a block:

```
output loki <loki_endpoint>
```

## Examples

```
{
    log default {
        output loki {env.LOKI_SERVER}
    }
}
```

Sends log output to a Loki server whose address is set as `LOKI_SERVER` in the environment Caddy is running in.

```
{
    log default {
        output loki https://myuser:password@loki.example.com {
           REGION iad  
        }
    }
}
```

Sends log messages to a Loki instance running at loki.example.com using basic authentication and with the `REGION=iad` label.