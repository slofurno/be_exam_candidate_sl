package main

import (
	"testing"
)

func TestValidPerson(t *testing.T) {
	invalidPeople := []Person{
		{"12345678", Name{"first", "middle", "last"}, "222-1111234"},
		{"1234567", Name{"first", "middle", "last"}, "222-111-1234"},
		{"12345678", Name{"", "middle", "last"}, "222-111-1234"},
		{"12345678", Name{"first", "middle", ""}, "222-111-1234"},
		{"12345678", Name{"averylongfirstname", "middle", "last"}, "2221-111-234"},
		{"12345678", Name{"first", "middle", "last"}, "a22-111-1234"},
		{"-12345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"123456789", Name{"first", "middle", "last"}, "222-111-1234"},
	}

	for _, p := range invalidPeople {
		if ok, _ := p.Valid(); ok {
			t.Error("should be invalid:", p)
		}

	}

	validPeople := []Person{
		{"12345678", Name{"first", "middle", "last"}, "222-111-1234"},
		{"12345678", Name{"first", "", "last"}, "222-111-1234"},
	}

	for _, p := range validPeople {
		if ok, _ := p.Valid(); !ok {
			t.Error("should be valid:", p)
		}

	}
}
