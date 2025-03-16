# piproxy

Reverse proxy for my rpi5

## Use

Create `services.json`. Example with a website server service:

```json
[
	{
		"name": "My website",
		"url": "http://localhost:8080",
		"endpoint": "/website"
	}
]
```

Run `cenv fix` and fill in the empty fields.

- `SERVICE_PATH` should be `services.json`
- `HOST` is either `localhost` for testing or the exposed rpi host name `<host>.local`.
- `PORT` is the port the proxy listens to

Start the services and run `go run cmd/main.go`.

## Roadmap

- [x] Set up proxy
- [ ] Run as daemon
- [ ] Let proxy start/stop services
- [ ] Simple dashboard for service diagnostics
- [ ] HTTPS/TLS
