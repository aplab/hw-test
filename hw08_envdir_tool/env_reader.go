package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"golang.org/x/exp/utf8string"
)

var (
	ErrUnableToReadDir  = errors.New("unable to read dir")
	ErrUnableToReadFile = errors.New("unable to read file")
	ErrUnableToOpenFile = errors.New("unable to open file")
	ErrNotADir          = errors.New("path is not a directory")
	ErrWrongVarName     = errors.New("wrong var name")
	ErrNotAFile         = errors.New("path is not a file")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	f, err := os.Stat(dir)
	if err != nil {
		return nil, ErrUnableToReadDir
	}
	if !f.IsDir() {
		return nil, ErrNotADir
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrUnableToReadDir
	}
	env := make(map[string]EnvValue)
	for _, f := range files {
		if f.IsDir() {
			return nil, ErrNotAFile
		}
		ev, err := readValueFromFile(dir, f.Name())
		if err != nil {
			return nil, err
		}
		env[f.Name()] = *ev
	}
	return env, nil
}

func readValueFromFile(dir, name string) (*EnvValue, error) {
	if !checkName(name) {
		return nil, ErrWrongVarName
	}
	file, err := os.Open(path.Join(dir, name))
	if err != nil {
		return nil, ErrUnableToOpenFile
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	bytes, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, ErrUnableToReadFile
	}
	return handleValue(bytes), nil
}

func handleValue(b []byte) *EnvValue {
	v := strings.ReplaceAll(strings.TrimRight(string(b), " \r\n\t"), "\x00", "\n")
	return &EnvValue{
		Value:      v,
		NeedRemove: len(v) == 0,
	}
}

func checkName(name string) bool {
	return !strings.ContainsRune(name, '=') && utf8string.NewString(name).IsASCII() && strings.TrimSpace(name) == name
}
