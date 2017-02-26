package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"strconv"
)

type invalidRecord struct {
	line    int
	message string
}

func (s *invalidRecord) ToRecord() []string {
	return []string{strconv.Itoa(s.line), s.message}
}

func processFile(name string) ([]*Person, []*invalidRecord, error) {
	f, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	var parsed []*Person
	var invalid []*invalidRecord
	line := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if line > 0 {
			if person, err := recordToPerson(record); person != nil {
				parsed = append(parsed, person)
			} else {
				invalid = append(invalid, &invalidRecord{line, err})
			}
		}
		line++
	}

	return parsed, invalid, nil
}

func main() {

	inputDir := flag.String("in", "input", "directory to watch for input")
	outputDir := flag.String("out", "output", "output directory")
	errorDir := flag.String("err", "errors", "error record directory")
	flag.Parse()

	store := newFileStore(*inputDir, *outputDir, *errorDir)
	job := newJob("example")
	job.processRecords(store)

	out, err := os.Create("output/example.json")

	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewEncoder(out).Encode(people); err != nil {
		log.Fatal(err)
	}

	out.Close()

	if len(invalid) > 0 {
		errorf, err := os.Create("errors/example.csv")
		defer errorf.Close()

		if err != nil {
			log.Fatal(err)
		}
		defer errorf.Close()
		writer := csv.NewWriter(errorf)
		for i := 0; i < len(invalid); i++ {
			_ = writer.Write(invalid[i].ToRecord())
		}

		writer.Flush()
		errorf.Write([]byte("LINE_NUM,ERROR_MSG\n"))
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add("./input")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			switch event.Op {
			case fsnotify.Create:
				log.Println("created file:", event.Name)
			default:
				log.Println(event.String())
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}

	}
}
