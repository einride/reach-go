package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/einride/reach-go/erb"
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
	sc := erb.NewClient(conn)
	for {
		if err := sc.Receive(context.Background()); err != nil {
			panic(err)
		}
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
}
