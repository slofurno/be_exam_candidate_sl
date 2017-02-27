package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"regexp"
)

var matchId = regexp.MustCompile(`(\S+)\.csv$`)

func extractId(p string) []string {
	_, filename := path.Split(p)
	return matchId.FindStringSubmatch(filename)
}

func main() {

	inputDir := flag.String("in", "input", "directory to watch for input")
	outputDir := flag.String("out", "output", "output directory")
	errorDir := flag.String("err", "errors", "error record directory")
	flag.Parse()

	store := newFileStore(*inputDir, *outputDir, *errorDir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(*inputDir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			switch event.Op {
			case fsnotify.Create:
				if id := extractId(event.Name); id != nil {
					job := newJob(id[1])
					log.Println("processing job id:", id[1])
					err := job.processRecords(store)
					if err != nil {
						log.Println("error processing records:", err)
						continue
					}
					err = job.writeResults(store)
					if err != nil {
						log.Println("error writing job output:", err)
						continue
					}
					err = job.writeErrors(store)
					if err != nil {
						log.Println("error writing job errors:", err)
						continue
					}

					if err := job.cleanup(store); err != nil {
						log.Println("error during job cleanup", err)
						continue
					}
				}
			default:
				log.Println(event.String())
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}

	}
}
