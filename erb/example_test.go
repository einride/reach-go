package erb_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"go.einride.tech/reach/erb"
)

func ExampleScanner() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// This example uses a mocked Reach unit that answers with recorded data.
	reach := newExampleReach(ctx)
	// Connect to the Emlid Reach Binary (ERB) protocol port of the Reach.
	conn, err := net.Dial("tcp", reach.Address)
	if err != nil {
		panic(err)
	}
	// Wrap the connection in an ERB protocol scanner.
	sc := erb.NewScanner(conn)
	// Scan 5 packets then exit.
	const maxPacketCount = 5
	packetCount := 0
	for sc.Scan() {
		// Handle packet.
		switch sc.ID() {
		case erb.IDVER:
			ver := sc.VER()
			fmt.Printf("VER:  %d.%d.%d\n", ver.High, ver.Medium, ver.Low)
		case erb.IDPOS:
			pos := sc.POS()
			fmt.Printf("POS:  %f,%f\n", pos.LatitudeDegrees, pos.LongitudeDegrees)
		case erb.IDSTAT:
			stat := sc.STAT()
			fmt.Printf("STAT: fix:%t type:%v satellites:%d\n", stat.HasFix, stat.FixType, stat.NumSVs)
		case erb.IDDOPS:
			dops := sc.DOPS()
			fmt.Printf("DOPS: pdop:%f\n", dops.Position)
		case erb.IDVEL:
			vel := sc.VEL()
			fmt.Printf("VEL:  speed:%dcm/s\n", vel.SpeedCentimetersPerSecond)
		case erb.IDSVI:
			svi := sc.SVI()
			fmt.Printf("SVI:  sattelites:%d\n", svi.NumSVs)
			for sc.ScanSVI() {
				sv := sc.SV()
				fmt.Printf("SV:   id:%v\n", sv.ID)
			}
		}
		packetCount++
		if packetCount >= maxPacketCount {
			break
		}
	}
	if sc.Err() != nil {
		panic(err)
	}
	if err := conn.Close(); err != nil {
		panic(err)
	}
	// Output:
	// VER:  0.1.0
	// POS:  57.777683,12.780540
	// STAT: fix:true type:Single satellites:20
	// DOPS: pdop:1.210000
	// VEL:  speed:0cm/s
}

type exampleReach struct {
	Address string
}

func newExampleReach(ctx context.Context) *exampleReach {
	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		if _, err := conn.Write(loadHexDump("testdata/hexdump.asta")); err != nil {
			panic(err)
		}
		if err := conn.Close(); err != nil {
			panic(err)
		}
		if err := lis.Close(); err != nil {
			panic(err)
		}
	}()
	return &exampleReach{
		Address: lis.Addr().String(),
	}
}

func loadHexDump(filename string) []byte {
	var data []byte
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 {
			continue
		}
		for _, field := range fields[1:] {
			b, err := strconv.ParseUint(field, 8, 8)
			if err != nil {
				panic(err)
			}
			data = append(data, byte(b))
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	return data
}
