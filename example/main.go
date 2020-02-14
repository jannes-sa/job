package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/satriarrrrr/job"
)

const (
	worker    = 5
	nmRoutine = "print-name"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

// Processor ...
type Processor struct {
}

// Validate ...
func (p Processor) Validate() (map[int]interface{}, bool) {
	// read file
	// put each line to map
	f, err := os.Open("in.txt")
	if err != nil {
		log.Println(">>> got error when validate:", err)
		return nil, false
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var (
		i     = 0
		tasks = make(map[int]interface{})
	)
	for fileScanner.Scan() {
		log.Println(">>> read line:", i)
		tasks[i] = fileScanner.Text()
		i++
	}

	return tasks, true
}

// Run ...
func (p Processor) Run(input job.ChanInputData) (interface{}, error) {
	s, ok := input.Data.(string)
	if !ok {
		fmt.Println(">>> non-string received:", input.Data)
		return input.Data, nil
	}

	fmt.Println(">>> string received: ", strings.Join(strings.Split(s, ","), "---"))
	return input.Data, nil
}

// Done ...
func (p Processor) Done(*job.OutputData) bool {
	return true
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	processor := Processor{}

	err := job.RunScheduler(worker, nmRoutine, processor)
	if err != nil {
		log.Println(">>> err: ", err)
	} else {
		log.Println(">>> success")
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
