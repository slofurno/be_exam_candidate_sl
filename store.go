package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

type bufferedReadCloser struct {
	buffer *bytes.Buffer
}

func (s bufferedReadCloser) Read(b []byte) (int, error) {
	return s.buffer.Read(b)
}

func (s bufferedReadCloser) Close() error {
	return nil
}

type bufferedWriteCloser struct {
	buffer *bytes.Buffer
}

func (s bufferedWriteCloser) Write(b []byte) (int, error) {
	return s.buffer.Write(b)
}

func (s bufferedWriteCloser) Close() error {
	return nil
}

type MemoryStore struct {
	input  bufferedReadCloser
	output bufferedWriteCloser
	errors bufferedWriteCloser
}

func (s *MemoryStore) OpenInput(id string) (io.ReadCloser, error) {
	return s.input, nil
}

func (s *MemoryStore) OpenOutput(id string) (io.WriteCloser, error) {
	return s.output, nil
}

func (s *MemoryStore) OpenError(id string) (io.WriteCloser, error) {
	return s.errors, nil
}

func newMemoryStory() *MemoryStore {
	return &MemoryStore{
		input:  bufferedReadCloser{&bytes.Buffer{}},
		output: bufferedWriteCloser{&bytes.Buffer{}},
		errors: bufferedWriteCloser{&bytes.Buffer{}},
	}
}
