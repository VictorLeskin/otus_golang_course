package hw07_file_copying

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTxtProgressBar struct {
	TxtProgressBar
}

func Test_TxtProgressBar_ctor(t *testing.T) {
	var t0 testTxtProgressBar

	assert.Equal(t, int64(0), t0.total)
	require.Equal(t, int64(0), t0.processed)

	t1 := NewTxtProgressBar(256, 50)

	assert.Equal(t, int64(256), t1.total)
	require.Equal(t, int64(0), t1.processed)
}

func Test_TxtProgressBar_Update(t *testing.T) {
	t1 := NewTxtProgressBar(256, 50)
	t1.Update(128)

	assert.Equal(t, int64(256), t1.total)
	require.Equal(t, int64(128), t1.processed)
}

func Test_TxtProgressBar_Render(t *testing.T) {
	t1 := NewTxtProgressBar(256, 20)

	t1.Render()
	//                  01234567890123456789
	assert.Equal(t, "\r[                    ] 0.0% (0/256)", t1.bar)

	t1.Update(128)
	t1.Render()
	//                  01234567890123456789
	assert.Equal(t, "\r[#########           ] 50.0% (128/256)", t1.bar)
}
