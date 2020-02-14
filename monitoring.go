package job

import "time"

func monitoring(
	worker int,
	nmRoutine string,
	logic logiclayer,
	input chan interface{},
	output chan correlated,
) (err error) {

	interval, err := time.ParseDuration("5s")
	if err != nil {
		return
	}

	for t := range time.Tick(interval) {
		print(t, nmRoutine, "TOTAL TASKS LEFT", routineStorage.numOfTask(nmRoutine), "STATUS", status.String(mappingStatusTasks[nmRoutine]))

		switch mappingStatusTasks[nmRoutine] {
		case restart:
			err = prepareRun(worker, nmRoutine, input, output)
			if err != nil {
				return
			}
		case done:
			return
		}
	}

	return
}
