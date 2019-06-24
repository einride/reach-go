package erb

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

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
			sc := NewScanner(bytes.NewReader(data))
			var buf bytes.Buffer
			for sc.Scan() {
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
			require.NoError(t, sc.Err())
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
