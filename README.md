# Reach Go

[![PkgGoDev][pkg-badge]][pkg]
[![GoReportCard][report-badge]][report]
[![Codecov][codecov-badge]][codecov]

[pkg-badge]: https://pkg.go.dev/badge/go.einride.tech/reach
[pkg]: https://pkg.go.dev/go.einride.tech/reach
[report-badge]: https://goreportcard.com/badge/go.einride.tech/reach
[report]: https://goreportcard.com/report/go.einride.tech/reach
[codecov-badge]: https://codecov.io/gh/einride/reach-go/branch/master/graph/badge.svg
[codecov]: https://codecov.io/gh/einride/reach-go

Go client for [Emlid Reach][emlid-reach] GNSS receivers.

[emlid-reach]: https://emlid.com/reach/

## Usage

```bash
$ go get -u go.einride.tech/reach
```

## Examples

### Reach ERB protocol data

```go
package main

import (
	"fmt"
	"net"

	"go.einride.tech/reach/erb"
)

func main() {
	// Connect to the Emlid Reach Binary (ERB) protocol port of the Reach.
	conn, err := net.Dial("tcp", "<REACH_ERB_ADDRESS>")
	if err != nil {
		panic(err)
	}
	// Wrap the connection in an ERB protocol scanner.
	sc := erb.NewScanner(conn)
	for sc.Scan() {
		// Handle packet.
		switch sc.ID() {
		case erb.IDVER:
			fmt.Printf("%v: %+v\n", sc.ID(), sc.VER())
		case erb.IDPOS:
			fmt.Printf("%v: %+v\n", sc.ID(), sc.POS())
		case erb.IDSTAT:
			fmt.Printf("%v: %+v\n", sc.ID(), sc.STAT())
		case erb.IDDOPS:
			fmt.Printf("%v: %+v\n", sc.ID(), sc.DOPS())
		case erb.IDVEL:
			fmt.Printf("%v: %+v\n", sc.ID(), sc.VEL())
		case erb.IDSVI:
			fmt.Printf("%v: %+v\n", sc.ID(), sc.SVI())
			for sc.ScanSVI() {
				fmt.Printf("%v: %+v\n", sc.ID(), sc.SV())
			}
		default:
			fmt.Printf("%v: %s\n", sc.ID(), hex.EncodeToString(sc.Bytes()))
		}
	}
	if sc.Err() != nil {
		panic(err)
	}
	if err := conn.Close(); err != nil {
		panic(err)
	}
}
```

_[Reference ≫][erb-protocol]_

[erb-protocol]: https://files.emlid.com/ERB.pdf
