package logic2

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"util/channel/lib/scheduler/job"
)

const (
	logicNm string = "logic2"
)

type logicSt struct {
	checkValidate int
	tasks         []Tasks
}

func (l logicSt) Validate() (bufferingTasks map[int]interface{}, state bool) {
	state = false

	if l.checkValidate == 2 {
		state = true

		bufferingTasks = make(map[int]interface{})
		for k, v := range l.tasks {
			bufferingTasks[k] = v
		}

		return
	}

	return
}

func (l logicSt) Run(receiverArg job.ChanInputData) (
	resp interface{},
	err error,
) {
	fmt.Println(time.Now(), logicNm, " => ", receiverArg.Data.(Tasks))

	failmap := map[int]bool{20: true, 30: true}
	if failmap[receiverArg.Data.(Tasks).task] {
		err = errors.New("FAILED SOMETHING POKOKNYA")
		return
	}

	resp = "RESPONSE " + strconv.Itoa(receiverArg.Data.(Tasks).task)

	return
}

// Done -
// return true => stop job and close all workers
// return false => will restart job start from validate -> run -> done again
func (l logicSt) Done(out *job.OutputData) (state bool) {
	fmt.Println(
		"RESULT DONE", (*out).Result, "\n",
		"TOTAL TASK", (*out).TotalTasks, "\n",
		"TOTAL TASK DONE", (*out).TotalTasksDone, "\n",
		"TOTAL TASK FAIL", (*out).TotalTasksFail, "\n",
		"TOTAL TASK PENDING", (*out).TotalTasksPending, "\n",
	)

	for _, v := range (*out).Err {
		fmt.Println(
			"ERROR =>", v.Err, "\n",
			"INPUT TASK =>", v.InputError, "\n",
		)
	}

	if (*out).TotalTasksFail == 0 {
		return true
	} else {
		return false
	}

	return
}

type Tasks struct {
	task       int
	taskString string
}

func RunScheduler() {

	var tasks []Tasks
	for i := 0; i < 5; i++ {
		tasks = append(tasks, Tasks{task: i, taskString: "YYYY"})
	}

	var l logicSt
	l.checkValidate = 2
	l.tasks = tasks

	err := job.RunScheduler(100, logicNm, &l)
	if err != nil {
		fmt.Println(err)
	}
}
