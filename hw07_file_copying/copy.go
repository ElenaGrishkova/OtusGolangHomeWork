package main

import (
	"errors"
	"io"
	"log"
	"os"

	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrDuplicateToPath       = errors.New("output path duplicates input path")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrDuplicateToPath
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer closeFile(fromFile)

	fromInfo, err := fromFile.Stat()
	if err != nil || !fromInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fromInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	println("limit before = ", limit)
	if limit == 0 || limit > fileSize {
		limit = fileSize
	}
	limit = min(limit, fileSize-offset)
	println("fileSize = ", fileSize)
	println("offset = ", offset)
	println("limit = ", limit)

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeFile(toFile)

	// set the offset
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	// create proxy reader
	barReader := bar.NewProxyReader(fromFile)
	// copy from proxy reader
	_, err = io.CopyN(toFile, barReader, limit)
	if err != nil {
		return err
	}
	// finish bar
	bar.Finish()

	return nil
}

func closeFile(fromFile *os.File) {
	err := fromFile.Close()
	if err != nil {
		log.Panicf("failed to close source file: %v", err)
	}
}
