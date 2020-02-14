package job

import "fmt"

type scheduler struct{}

func (s scheduler) run(
	routine int,
	nmRoutine string,
	tasks map[int]interface{},
	input chan interface{},
	output chan correlated,
) {
	routineStorage.setTasks(nmRoutine, tasks)

	if mappingStatusTasks[nmRoutine] == restart {
		input, output = make(chan interface{}), make(chan correlated)
	}

	mappingStatusTasks[nmRoutine] = running
	for i := 0; i < routine; i++ {
		go worker(input, output, nmRoutine)
	}

	go sendinput(nmRoutine, input)
	getOutput(routineStorage.numOfTask(nmRoutine), nmRoutine, output)

	if mappingStatusTasks[nmRoutine] == done || mappingStatusTasks[nmRoutine] == restart {
		fmt.Println("######### TEAR DOWN CHANNEL AND GO ROUTINE #########")
		tearDown(input, output)
	}
}
