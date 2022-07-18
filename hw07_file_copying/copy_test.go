package main

import (
	"fmt"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/aplab/hw-test/hw07_file_copying/progressbar"
	"github.com/hlubek/readercomp"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("full copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out1.txt", 0, 0)
		require.NoError(t, err)
		result, err := readercomp.FilesEqual("out1.txt", "testdata/out_offset0_limit0.txt")
		require.True(t, result)
		require.NoError(t, err)
		err = os.Remove("out1.txt")
		require.NoError(t, err)
	})

	t.Run("copy limit 10", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out2.txt", 0, 10)
		require.NoError(t, err)
		result, err := readercomp.FilesEqual("out2.txt", "testdata/out_offset0_limit10.txt")
		require.True(t, result)
		require.NoError(t, err)
		err = os.Remove("out2.txt")
		require.NoError(t, err)
	})

	t.Run("copy limit 1000", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out3.txt", 0, 1000)
		require.NoError(t, err)
		result, err := readercomp.FilesEqual("out3.txt", "testdata/out_offset0_limit1000.txt")
		require.True(t, result)
		require.NoError(t, err)
		err = os.Remove("out3.txt")
		require.NoError(t, err)
	})

	t.Run("copy limit 10000", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out4.txt", 0, 10000)
		require.NoError(t, err)
		result, err := readercomp.FilesEqual("out4.txt", "testdata/out_offset0_limit10000.txt")
		require.True(t, result)
		require.NoError(t, err)
		err = os.Remove("out4.txt")
		require.NoError(t, err)
	})

	t.Run("copy offset 100 limit 1000", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out5.txt", 100, 1000)
		require.NoError(t, err)
		result, err := readercomp.FilesEqual("out5.txt", "testdata/out_offset100_limit1000.txt")
		require.True(t, result)
		require.NoError(t, err)
		err = os.Remove("out5.txt")
		require.NoError(t, err)
	})

	t.Run("copy offset 6000 limit 1000", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out6.txt", 6000, 1000)
		require.NoError(t, err)
		result, err := readercomp.FilesEqual("out6.txt", "testdata/out_offset6000_limit1000.txt")
		require.True(t, result)
		require.NoError(t, err)
		err = os.Remove("out6.txt")
		require.NoError(t, err)
	})

	t.Run("src file not exists", func(t *testing.T) {
		err := Copy("testdata/unknown.txt", "out7.txt", 0, 0)
		require.ErrorIs(t, err, ErrSrcFileNotFound)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out7.txt", 1<<20, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("unable to write", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/etc/i-do-not-have-permission-to-write-there.txt", 0, 0)
		fmt.Println(err)
		require.ErrorIs(t, err, ErrUnableToCreateDstFile)
	})

	t.Run("progress", func(t *testing.T) {
		p := progressbar.NewProgressbar(1000)
		p.SetValue(100)
		require.Equal(t, 10, p.GetPercentage())
	})

	t.Run("progress value overflow", func(t *testing.T) {
		p := progressbar.NewProgressbar(1000)
		p.SetValue(10000)
		require.Equal(t, 100, p.GetPercentage())
	})

	t.Run("progress length", func(t *testing.T) {
		p := progressbar.NewProgressbar(1000)
		p.SetValue(-12345)
		require.Equal(t, 105, utf8.RuneCountInString(p.String()))

		p.SetValue(100)
		require.Equal(t, 106, utf8.RuneCountInString(p.String()))

		p.SetValue(1000)
		require.Equal(t, 107, utf8.RuneCountInString(p.String()))

		p.Finish()
		require.Equal(t, 107, utf8.RuneCountInString(p.String()))
	})
}
