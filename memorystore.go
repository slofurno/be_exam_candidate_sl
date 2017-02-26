package main

import (
	"bytes"
	"io"
)

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
