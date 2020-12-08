package erb

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestScanner_HexDump(t *testing.T) {
	for _, tt := range []struct {
		inputFile  string
		goldenFile string
	}{
		{inputFile: "testdata/hexdump.empty", goldenFile: "testdata/hexdump.empty.golden"},
		{inputFile: "testdata/hexdump.asta", goldenFile: "testdata/hexdump.asta.golden"},
	} {
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			data := loadHexDump(t, tt.inputFile)
			sc := NewClient(newMockConn(bytes.NewReader(data)))
			var buf bytes.Buffer
			for {
				if err := sc.Receive(context.Background()); err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					require.NoError(t, err)
				}
				switch sc.ID() {
				case IDVER:
					_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.VER())
				case IDPOS:
					_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.POS())
				case IDSTAT:
					_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.STAT())
				case IDDOPS:
					_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.DOPS())
				case IDVEL:
					_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.VEL())
				case IDSVI:
					_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.SVI())
					for sc.ScanSVI() {
						_, _ = fmt.Fprintf(&buf, "%v: %+v\n", sc.ID(), sc.SV())
					}
				default:
					_, _ = fmt.Fprintf(&buf, "%v: %s\n", sc.ID(), hex.EncodeToString(sc.Bytes()))
				}
			}
			if shouldUpdateGoldenFiles() {
				require.NoError(t, ioutil.WriteFile(tt.goldenFile, buf.Bytes(), 0644))
			}
			requireGoldenFileContent(t, tt.goldenFile, buf.String())
		})
	}
}

func loadHexDump(t *testing.T, filename string) []byte {
	t.Helper()
	var data []byte
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, f.Close())
	}()
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 {
			continue
		}
		for _, field := range fields[1:] {
			b, err := strconv.ParseUint(field, 8, 8)
			require.NoError(t, err)
			data = append(data, byte(b))
		}
	}
	require.NoError(t, sc.Err())
	return data
}

type mockConn struct {
	r io.Reader
}

func newMockConn(r io.Reader) *mockConn {
	return &mockConn{r: r}
}

func (m mockConn) Read(p []byte) (n int, err error) {
	return m.r.Read(p)
}

func (m mockConn) Close() error {
	return nil
}

func (m mockConn) SetReadDeadline(time.Time) error {
	return nil
}
