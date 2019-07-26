package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

var fileDir = "./data"
var dict map[string]int
var mapMutex = sync.RWMutex{}

//Job ...
type Job struct {
	id       int
	fileName string
}

//MaxWorker ...
const MaxWorker int = 10

var jobs = make(chan Job, MaxWorker)

func countWordsFromFile(fileName string) {
	st, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", fileDir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(string(st))

	for i := 0; i < len(words); i++ {
		mapMutex.Lock()
		key := strings.ToLower(words[i])
		count := dict[key]
		dict[key] = count + 1
		mapMutex.Unlock()
	}

	fmt.Println("end of count file ", fileName)
}

func printCounter(done chan bool) {
	for key, value := range dict {
		fmt.Println(key, " = ", value)
	}

	done <- true
}

func createWorkerPool(numberOfWorker int) {
	var wg sync.WaitGroup

	for i := 0; i < numberOfWorker; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		countWordsFromFile(job.fileName)
	}

	wg.Done()
}

func getAllJobs() {
	files, err := ioutil.ReadDir(fileDir)
	if err != nil {
		log.Fatal(err)
	}

	for i, f := range files {
		job := Job{i, f.Name()}
		jobs <- job
	}

	close(jobs)
}

func main() {
	dict = make(map[string]int)

	go getAllJobs()
	createWorkerPool(MaxWorker)

	done := make(chan bool)
	go printCounter(done)
	<-done

	time.Sleep(1 * time.Second)
	fmt.Println("\nend main")
}
