package main

import (
	"io"
	"os"
	"path"
)

type Store interface {
	OpenInput(id string) (io.ReadCloser, error)
	OpenOutput(id string) (io.WriteCloser, error)
	OpenError(id string) (io.WriteCloser, error)
	RemoveInput(id string) error
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
	return os.Open(path.Join(s.inputDir, id+".csv"))
}

func (s *FileStore) OpenOutput(id string) (io.WriteCloser, error) {
	return os.Create(path.Join(s.outputDir, id+".json"))
}

func (s *FileStore) OpenError(id string) (io.WriteCloser, error) {
	return os.Create(path.Join(s.errorDir, id+".csv"))
}

func (s *FileStore) RemoveInput(id string) error {
	return os.Remove(path.Join(s.inputDir, id+".csv"))
}
