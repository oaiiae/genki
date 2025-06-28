# 元気 Genki

Genki is a japanese word used when asking how well a person is.

This repository contains somewhat optimized yet simple helpers for creating healthcheck probes to use in Go HTTP handlers.

## Usage example

```go
func _() {
	db, err := sql.Open("postgres", "...")
	if err != nil {
		// ...
	}

	http.Handle("/healthz", genki.Checks{
		"db": db.PingContext,
	}.Handler(func(err error) { log.Print(err) }))
}
```
