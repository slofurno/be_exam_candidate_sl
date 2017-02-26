package main

import (
	"encoding/csv"
)

type Job struct {
	id      string
	results []*Person
	errors  []*invalidRecord
}

func newJob(id string) *Job {
	return &Job{id: id}
}

func (s *Job) writeResults(store Store) error {
	output, err := store.OpenOutput(s.id)

	if err != nil {
		return err
	}

	if err = json.NewEncoder(output).Encode(s.results); err != nil {
		return err
	}

	if err = output.Close(); err != nil {
		return err
	}

	if len(s.errors) > 0 {
		errorOut, err := store.OpenError(s.id)
		if err != nil {
			return err
		}

		errorOut.Write([]byte("LINE_NUM,ERROR_MSG\n"))
		writer := csv.NewWriter(errorOut)
		for i := 0; i < len(s.errors); i++ {
			_ = writer.Write(s.errors[i].ToRecord())
		}
		writer.Flush()
		if err = errorOut.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Job) processRecords(store Store) error {
	input, err := store.OpenInput(s.id)

	if err != nil {
		return err
	}

	defer input.Close()
	reader := csv.NewReader(f)
	line := -1

	for {
		line++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			s.errors = append(s.errors, &invalidRecord{line, "invalid csv record"})
			continue
		}

		if line > 0 {
			if person, err := recordToPerson(record); person != nil {
				s.results = append(s.results, person)
			} else {
				s.errors = append(s.errors, &invalidRecord{line, err})
			}
		}
	}

	return nil
}
