package placement

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeKey(t *testing.T) {
	re := require.New(t)
	for i := 0; i < 256; i++ {
		nid := ((((0x01<<8 + int64(i)) << 16) + 0) << 32) + 1
		nid2 := 0x01_0000_00_0000_0000 | int64(i)<<48 | 0<<32 | 1
		re.Equal(nid, nid2)
		id := (nid & 0x00_FFFF_00_0000_0000) >> 48
		re.Equal(int64(i), id)
	}
}
