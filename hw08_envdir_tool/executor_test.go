package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("executor", func(t *testing.T) {
		pwd, _ := os.Getwd()
		code := RunCmd([]string{
			"/usr/bin/cat",
			path.Join(pwd, "testdata/env/BAR"),
		}, Environment{
			"BAR": EnvValue{
				Value:      "BAR",
				NeedRemove: false,
			},
		})
		require.Equal(t, code, 0)
	})

	t.Run("executor", func(t *testing.T) {
		pwd, _ := os.Getwd()
		code := RunCmd([]string{
			"/bin/bash",
			"unknowncommand",
			path.Join(pwd, "testdata/env/BAR"),
		}, Environment{
			"FOO": EnvValue{
				Value:      "BAR",
				NeedRemove: false,
			},
		})
		require.NotEqual(t, 0, code)
	})

	t.Run("set var", func(t *testing.T) {
		err := setVariable("FOO1", EnvValue{
			Value:      "BAR1",
			NeedRemove: false,
		})
		require.NoError(t, err)
		require.Equal(t, os.Getenv("FOO1"), "BAR1")
	})

	t.Run("unset var error", func(t *testing.T) {
		err := setVariable("FO=O1", EnvValue{
			Value:      "BAR1",
			NeedRemove: false,
		})
		require.ErrorIs(t, err, ErrUnableToSetVar)
	})

	t.Run("unset var", func(t *testing.T) {
		err := setVariable("FOO1", EnvValue{
			Value:      "BAR1",
			NeedRemove: true,
		})
		require.NoError(t, err)
		require.Equal(t, "", os.Getenv("FOO1"))
	})
}
