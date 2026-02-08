package hw10programoptimization

import (
	"archive/zip"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Benchmark_GetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()

	require.Equal(b, 1, len(r.File))

	data, err := r.File[0].Open()
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		func() {
			_, _ = originalGetDomainStat(data, "biz")
		}()
	}
}

// go test -bench=Benchmark_GetDomainStat .
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

func Test_CmpOriginalAndMy(t *testing.T) {
	benchOriginal := func(b *testing.B) {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(t, err)

		b.StartTimer()
		stat, err := originalGetDomainStat(data, "biz")
		b.StopTimer()
		require.NoError(t, err)
		require.True(t, len(stat) > 0)
	}

	benchMy := func(b *testing.B) {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(t, err)

		b.StartTimer()
		stat, err := GetDomainStat(data, "biz")
		b.StopTimer()
		require.NoError(t, err)
		require.True(t, len(stat) > 0)
	}

	resultOriginal := testing.Benchmark(benchOriginal)
	result := testing.Benchmark(benchMy)
	fmt.Printf("time used: %s\n", resultOriginal.T)
	fmt.Printf("time used: %s\n", result.T)
}
