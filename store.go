package main

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Store interface {
	OpenInput(id string) (io.ReadCloser, error)
	OpenOutput(id string) (io.WriteCloser, error)
	OpenError(id string) (io.WriteCloser, error)
}

type FileStore struct {
	inputDir  string
	outputDir string
	errorDir  string
}

func newFileStore(inputDir, outputDir, errorDir string) *FileStore {
	os.Mkdir(inputDir, os.ModePerm)
	os.Mkdir(outputDir, os.ModePerm)
	os.Mkdir(errorDir, os.ModePerm)

	return &FileStore{inputDir, outputDir, errorDir}
}

func (s *FileStore) OpenInput(id string) (io.ReadCloser, error) {
	return os.Open(s.inputDir + "/" + id + ".csv")
}

func (s *FileStore) OpenOutput(id string) (io.WriteCloser, error) {
	return os.Create(fmt.Sprintf("%s/%s.json", s.outputDir, id))
}

func (s *FileStore) OpenError(id string) (io.WriteCloser, error) {
	return os.Create(fmt.Sprintf("%s/%s.csv", s.errorDir, id))
}
