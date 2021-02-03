package bytesutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBytesUtil(t *testing.T) {
	require.Equal(t, ByteSize(1024), "1K")
	require.Equal(t, ByteSize(1024*1024), "1M")
	require.Equal(t, ByteSize(1000), "1000B")
	require.Equal(t, ByteSize(0), "0B")

	var (
		sizeStr string
		size    uint64
		err     error
	)

	sizeStr = "0B"
	size, err = ToBytes(sizeStr)
	require.NoError(t, err)
	require.Equal(t, size, uint64(0))

	sizeStr = "1M"
	size, err = ToBytes(sizeStr)
	require.NoError(t, err)
	require.Equal(t, size, uint64(1024*1024))
}
