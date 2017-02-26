package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
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
	validId = regexp.MustCompile(`^\d{8}$`)
)

func (p *Person) Valid() (bool, string) {
	if id, _ := validId.MatchString(p.Id); !id {
		return false, "invalid id"
	}

	if first, _ := regexp.MatchString(`^.{1,15}$`, p.Name.First); !first {
		return false, "invalid first name"
	}

	if middle, _ := regexp.MatchString(`^.{0,15}$`, p.Name.Middle); !middle {
		return false, "invalid middle name"
	}

	if last, _ := regexp.MatchString(`^.{1,15}$`, p.Name.Last); !last {
		return false, "invalid last name"
	}

	if phone, _ := regexp.MatchString(`^\d{3}-\d{3}-\d{4}$`, p.Phone); !phone {
		return false, "invalid phone number"
	}

	return true, ""
}

func fromRecord(record []string) (*Person, string) {
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
			if person, err := fromRecord(record); person != nil {
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
	name := "example.csv"
	os.Mkdir("output", os.ModePerm)
	os.Mkdir("errors", os.ModePerm)

	people, invalid, err := processFile(name)

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
		_ = writer.Flush()
		errorf.Write([]byte("LINE_NUM,ERROR_MSG\n"))
	}

	fmt.Println("vim-go")
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
