package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("check name", func(t *testing.T) {
		result := checkName("asd")
		require.True(t, result)

		result = checkName("VA=R")
		require.False(t, result)

		result = checkName("VAЯR")
		require.False(t, result)

		result = checkName(" VAR")
		require.False(t, result)
	})

	t.Run("check handle value", func(t *testing.T) {
		result := handleValue([]byte(""))
		require.True(t, result.NeedRemove)

		result = handleValue([]byte("hello"))
		require.False(t, result.NeedRemove)

		result = handleValue([]byte("hel\x00lo"))
		require.Equal(t, &EnvValue{
			Value:      "hel\nlo",
			NeedRemove: false,
		}, result)

		result = handleValue([]byte("  hello  \t  "))
		require.Equal(t, &EnvValue{
			Value:      "  hello",
			NeedRemove: false,
		}, result)
	})

	t.Run("read value from file", func(t *testing.T) {
		_, err := readValueFromFile("/dev", "nul=l")
		require.ErrorIs(t, err, ErrWrongVarName)

		_, err = readValueFromFile("/etc", "nonexistentfile123")
		require.ErrorIs(t, err, ErrUnableToOpenFile)

		result, err := readValueFromFile("testdata/env", "BAR")
		require.NoError(t, err)
		require.Equal(t, EnvValue{
			Value:      "bar",
			NeedRemove: false,
		}, *result)

		result, err = readValueFromFile("testdata/env", "EMPTY")
		require.NoError(t, err)
		require.Equal(t, EnvValue{
			Value:      "",
			NeedRemove: true,
		}, *result)

		result, err = readValueFromFile("testdata/env", "FOO")
		require.NoError(t, err)
		require.Equal(t, EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		}, *result)

		result, err = readValueFromFile("testdata/env", "HELLO")
		require.NoError(t, err)
		require.Equal(t, EnvValue{
			Value:      "\"hello\"",
			NeedRemove: false,
		}, *result)

		result, err = readValueFromFile("testdata/env", "UNSET")
		require.NoError(t, err)
		require.Equal(t, EnvValue{
			Value:      "",
			NeedRemove: true,
		}, *result)
	})

	t.Run("read dir", func(t *testing.T) {
		_, err := ReadDir("/notfounddir")
		require.ErrorIs(t, err, ErrUnableToReadDir)

		_, err = ReadDir("testdata/env/UNSET")
		require.ErrorIs(t, err, ErrNotADir)

		_, err = ReadDir("testdata/")
		require.ErrorIs(t, err, ErrNotAFile)

		_, err = ReadDir("/root")
		require.ErrorIs(t, err, ErrUnableToReadDir)

		result, _ := ReadDir("testdata/env")
		require.Equal(t, len(result), 5)
	})
}
