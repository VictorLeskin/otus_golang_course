package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

var emails_TestDomainStat []string = []string{
	"aliquid_qui_ea@Browsedrive.gov",
	"mLynch@broWsecat.com",
	"RoseSmith@Browsecat.com",
	"5Moore@Teklist.net",
	"nulla@Linktype.com",
}

func Benchmark_GetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()

	require.Equal(b, 1, len(r.File))

	data, err := r.File[0].Open()
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		func() {
			_, _ = original_GetDomainStat(data, "biz")
		}()
	}
}

func Benchmark_GetDomainStatMy(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()

	require.Equal(b, 1, len(r.File))

	data, err := r.File[0].Open()
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		func() {
			_, _ = GetDomainStat(data, "biz")
		}()
	}
}
