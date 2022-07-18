package main

import (
	"io"
	"os"
	"syscall"

	"github.com/aplab/hw-test/hw07_file_copying/progressbar"
	"github.com/pkg/errors"
)

const (
	DefaultBufferSize int64 = 4 * 1 << 10
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSrcFileNotFound       = errors.New("source file not found")
	ErrDstFileAlreadyExists  = errors.New("destination file already exists")
	ErrUnableToOpenSrcFile   = errors.New("unable to open source file")
	ErrUnableToOpenDstFile   = errors.New("unable to open destination file")
	ErrUnableToCreateDstFile = errors.New("unable to create destination file")
	ErrWrongDestinationFile  = errors.New("wrong destination file")
	ErrUnableToSetOffset     = errors.New("unable to set offset")
	ErrRead                  = errors.New("read error")
	ErrWrite                 = errors.New("write error")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcStat, err := getSrcStat(fromPath)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnableToOpenSrcFile
	}
	defer srcFile.Close()
	if offset > srcStat.Size() {
		return ErrOffsetExceedsFileSize
	}
	err = checkDstPath(toPath)
	if err != nil {
		return err
	}
	dstFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcStat.Mode())
	if err != nil {
		return ErrUnableToCreateDstFile
	}
	defer dstFile.Close()
	dstStat, err := dstFile.Stat()
	if err != nil {
		return ErrUnableToOpenDstFile
	}
	bufferSize := DefaultBufferSize
	sys := dstStat.Sys()
	if statT, ok := sys.(syscall.Stat_t); ok {
		bufferSize = statT.Blksize
	}
	buf := make([]byte, bufferSize)
	newOffset, err := srcFile.Seek(offset, 0)
	if err != nil || newOffset != offset {
		return ErrUnableToSetOffset
	}
	if limit == 0 {
		limit = srcStat.Size()
	}
	limitBytes := srcStat.Size() - newOffset
	if limit < limitBytes {
		limitBytes = limit
	}
	pb := progressbar.NewProgressbar(limitBytes)
	pb.Print()
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return ErrRead
		}
		if n == 0 {
			break
		}
		if limitBytes > int64(n) {
			limitBytes -= int64(n)
		} else {
			n = int(limitBytes)
			limitBytes -= int64(n)
		}
		if _, err := dstFile.Write(buf[:n]); err != nil {
			return ErrWrite
		}
		if limitBytes == 0 {
			break
		}
		pbOldPrc := pb.GetPercentage()
		pb.SetValue(pb.GetLimit() - limitBytes)
		if pbOldPrc != pb.GetPercentage() {
			pb.Print()
		}
	}
	pb.Finish()
	pb.Print()
	return nil
}

func getSrcStat(path string) (os.FileInfo, error) {
	srcStat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrSrcFileNotFound
		}
		return nil, ErrUnsupportedFile
	}
	if srcStat.IsDir() {
		return nil, ErrUnsupportedFile
	}
	return srcStat, nil
}

func checkDstPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		// unable to overwrite file to prevent third party file overwrite attack
		return ErrDstFileAlreadyExists
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ErrWrongDestinationFile
	}
	return nil
}
