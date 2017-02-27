package main

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

var matchId = regexp.MustCompile(`(\S+)\.csv$`)

func extractId(p string) []string {
	_, filename := filepath.Split(p)
	return matchId.FindStringSubmatch(filename)
}

type logger interface {
	Printf(string, ...interface{})
	Println(...interface{})
	Fatal(...interface{})
}

type logLogger struct{}

func (s logLogger) Printf(f string, a ...interface{}) {
	log.Printf(f, a...)
}

func (s logLogger) Println(a ...interface{}) {
	log.Println(a...)
}

func (s logLogger) Fatal(a ...interface{}) {
	log.Fatal(a...)
}

type recordProcessor struct {
	store Store
	logger
}

func (s *recordProcessor) processFile(filename string) {
	if id := extractId(filename); id != nil {
		job := newJob(id[1])
		s.logger.Println("processing job id:", id[1])

		if err := job.processRecords(s.store); err != nil {
			s.logger.Println("error processing records:", err)
			return
		}
		if err := job.writeResults(s.store); err != nil {
			s.logger.Println("error writing job output:", err)
			return
		}
		if err := job.writeErrors(s.store); err != nil {
			s.logger.Println("error writing job errors:", err)
			return
		}
		if err := job.cleanup(s.store); err != nil {
			s.logger.Println("error during job cleanup", err)
			return
		}
	}
}

func (s *recordProcessor) Run(inputDir string) {

	files, _ := ioutil.ReadDir(inputDir)

	for _, file := range files {
		if !file.IsDir() {
			s.processFile(file.Name())
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.logger.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(inputDir)
	if err != nil {
		s.logger.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			switch event.Op {
			case fsnotify.Create:
				s.processFile(event.Name)
			default:
				s.logger.Println(event.String())
			}
		case err := <-watcher.Errors:
			s.logger.Println("error:", err)
		}

	}
}
