package main

import (
	"regexp"
)

type Name struct {
	First  string `json:"first"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last"`
}

type Person struct {
	Id    string `json:"id"`
	Name  Name   `json:"name"`
	Phone string `json:"phone"`
}

var (
	validId     = regexp.MustCompile(`^\d{8}$`)
	validFirst  = regexp.MustCompile(`^.{1,15}$`)
	validMiddle = regexp.MustCompile(`^.{0,15}$`)
	validLast   = regexp.MustCompile(`^.{1,15}$`)
	validPhone  = regexp.MustCompile(`^\d{3}-\d{3}-\d{4}$`)
)

func (p *Person) Valid() (bool, string) {
	if !validId.MatchString(p.Id) {
		return false, "invalid id"
	}

	if !validFirst.MatchString(p.Name.First) {
		return false, "invalid first name"
	}

	if !validMiddle.MatchString(p.Name.Middle) {
		return false, "invalid middle name"
	}

	if !validLast.MatchString(p.Name.Last) {
		return false, "invalid last name"
	}

	if !validPhone.MatchString(p.Phone) {
		return false, "invalid phone number"
	}

	return true, ""
}

func recordToPerson(record []string) (*Person, string) {
	if len(record) != 5 {
		return nil, "invalid record length"
	}

	p := &Person{
		Id: record[0],
		Name: Name{
			First:  record[1],
			Middle: record[2],
			Last:   record[3],
		},
		Phone: record[4],
	}

	if ok, err := p.Valid(); !ok {
		return nil, err
	}

	return p, ""
}
