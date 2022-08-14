package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()
	require.Equal(b, 1, len(r.File))
	data, err := r.File[0].Open()
	require.NoError(b, err)
	_, err = GetDomainStat(data, "biz")
	require.NoError(b, err)
}
