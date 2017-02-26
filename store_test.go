package main

import (
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
