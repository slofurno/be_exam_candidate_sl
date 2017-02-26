package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestStoreWrite(t *testing.T) {
	people := []*Person{
		{"12345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"22345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"32345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"42345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"52345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"62345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"72345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"82345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"92345678", Name{"first", "middle", "last"}, "222-111-1234"},
	}

	job := newJob("test")
	store := newMemoryStory()
	job.results = people

	if err := job.writeResults(store); err != nil {
		t.Error(err)
	}

	var written []*Person
	if err := json.NewDecoder(store.output.buffer).Decode(&written); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(people, written) {
		t.Error("invalid json output written")
	}
}

func TestErrorOutput(t *testing.T) {

	expected := []byte(`LINE_NUM,ERROR_MSG
4,4
5,5
6,6
7,7
8,8
`)

	errors := []*invalidRecord{
		{4, "4"},
		{5, "5"},
		{6, "6"},
		{7, "7"},
		{8, "8"},
	}

	job := newJob("test")
	job.errors = errors
	store := newMemoryStory()

	if err := job.writeErrors(store); err != nil {
		t.Error(err)
	}

	written := store.errors.buffer.Bytes()

	if bytes.Compare(expected, written) != 0 {
		t.Error("unexpected error output")
	}
}

func TestReadInput(t *testing.T) {
	input := []byte(`
INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,steve,a,lofurno,267-663-9604
87654321,esteban,m,something,123-123-1234
8765432,esteban,m,something,123-123-1234
87654321,esteban,,something,123-123-1234
`)

	job := newJob("test")
	store := newMemoryStory()

	store.input.buffer = bytes.NewBuffer(input)

	if err := job.processRecords(store); err != nil {
		t.Error(err)
	}

	if len(job.results) != 3 {
		t.Errorf("expected 3 valid records read, got %d\n", len(job.results))
	}

	if len(job.errors) != 1 {
		t.Errorf("expected 1 invalid records read, got %d\n", len(job.errors))
	}
}
