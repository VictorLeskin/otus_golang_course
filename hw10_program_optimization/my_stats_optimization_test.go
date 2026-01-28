//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	mb          uint64 = 1 << 20
	memoryLimit uint64 = 30 * mb

	timeLimit = 300 * time.Millisecond
)


func BenchmarkGetDomainStat_Time_And_Memory(t *testing.B) {
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

	require.Equal(t, expectedBizStat, stat)

	result := testing.Benchmark(bench)
	mem := result.MemBytes
	t.Logf("time used: %s / %s", result.T, timeLimit)
	t.Logf("memory used: %dMb / %dMb", mem/mb, memoryLimit/mb)

	require.Less(t, int64(result.T), int64(timeLimit), "the program is too slow")
	require.Less(t, mem, memoryLimit, "the program is too greedy")
}

