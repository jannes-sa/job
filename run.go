package job

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	mappingTasks       = make(map[string]map[int]interface{})
	mappingStatusTasks = make(map[string]status)
	debug              = true
)

type status uint8

const (
	preparing = iota
	running
	restart
	done
)

func (s status) String() string {
	str := map[status]string{
		preparing: "preparing",
		running:   "running",
		restart:   "restart",
		done:      "done",
	}

	v, ok := str[s]
	if !ok {
		return "UNDEFINED"
	}
	return v
}

// RunScheduler - Running Scheduler
func RunScheduler(
	worker int,
	nmRoutine string,
	logic logiclayer,
) (err error) {

	input, output := make(chan interface{}), make(chan correlated)

	err = registerLogic(nmRoutine, logic, input, output)
	if err != nil {
		return
	}

	mappingStatusTasks[nmRoutine] = preparing
	err = prepareRun(worker, nmRoutine, input, output)
	if err != nil {
		return
	}

	monitoring(worker, nmRoutine, logic, input, output)
	return
}

func prepareRun(
	worker int,
	nmRoutine string,
	input chan interface{},
	output chan correlated,
) (err error) {
	var msg string
	tasks, state := logicRun[nmRoutine].Validate()
	if !state {
		msg = "VALIDATE JOB" + nmRoutine + "FALSE"
		println(msg)
		err = errors.New(msg)
		return
	}

	if len(tasks) == 0 {
		msg = "VALIDATE JOB" + nmRoutine + "TASKS" + strconv.Itoa(len(tasks))
		println(msg)
		err = errors.New(msg)
		return
	}

	var sch scheduler
	sch.run(worker, nmRoutine, tasks, input, output)
	return
}

func print(msg ...interface{}) {
	if debug {
		fmt.Println(msg)
	}
}
