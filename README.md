## autocert-reverseproxy

This tool acts as reverse proxy in front of your service, and can request TLS certs on your behalf

Usage:

```bash

$ autocert-reverseproxy -h
golang reverse proxy that handles cert provisioning as well

Usage:
  autocert-reverseproxy [flags]

Flags:
      --allowed-host stringArray   hostnames allowed to request cert
      --cert-cache-dir string      dir to cache certs to
  -h, --help                       help for autocert-reverseproxy
      --http-addr string           http address (default ":80")
      --https-addr string          https address (default ":443")
      --upstream string            upstream service (default "http://localhost:8080") 

```