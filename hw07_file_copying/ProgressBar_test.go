package hw07_file_copying

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testProgressBar struct {
	ProgressBar
}

func Test_ProgressBar_ctor(t *testing.T) {
	var t0 testProgressBar

	assert.Equal(t, int64(0), t0.expected)
	require.Equal(t, int64(0), t0.processed)

	t1 := NewProgressBar(256)

	assert.Equal(t, int64(256), t1.expected)
	require.Equal(t, int64(0), t1.processed)
}
