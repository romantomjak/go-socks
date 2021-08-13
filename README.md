# go-socks

[![Build Status](https://travis-ci.org/r00m/go-socks.svg?branch=master)](https://travis-ci.org/r00m/go-socks)
[![Coverage Status](https://coveralls.io/repos/github/r00m/go-socks/badge.svg?branch=master)](https://coveralls.io/github/r00m/go-socks?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/r00m/go-socks)](https://goreportcard.com/report/github.com/r00m/go-socks)
[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://github.com/r00m/go-socks/blob/master/LICENSE)

SOCKS (SOCKS4, SOCKS4a, SOCKS5) proxy library for Go

---

## Requirements

- Go 1.8

## Client example

```go
package main

import (
	"fmt"
	"io"

	"github.com/romantomjak/go-socks"
)

func main() {
	// 1. connect to a SOCKS server
	client, err := socks.NewV4Client("127.0.0.1:1234")
	if err != nil {
		panic(err)
	}

	// 2. instruct the server to relay connections to golang.org
	conn, err := client.Connect("142.250.187.241:80", "roman")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 3. fetch index.html via the relayed connection
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	body, err := io.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", body)
}
```

## Testing

```
$ go test
```

## License

MIT