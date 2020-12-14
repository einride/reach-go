package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"go.einride.tech/reach/erb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: reachctl <host:port>")
		os.Exit(1)
	}
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	sc := erb.NewScanner(conn)
	for sc.Scan() {
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
		panic(sc.Err())
	}
}
