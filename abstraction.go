package job

import "fmt"

var (
	logicRun = make(map[string]logiclayer)
)

type logiclayer interface {
	Validate() (tasks map[int]interface{}, state bool)
	Run(ChanInputData) (interface{}, error)
	Done(*OutputData) bool
}

// RegisterLogic - Register Logic Inside Scheduler
func registerLogic(
	nmRoutine string,
	logic logiclayer,
	input chan interface{},
	output chan correlated,
) (err error) {
	if _, ok := logicRun[nmRoutine]; ok {
		msg := "failed Registered Logic " + nmRoutine + "Already Registered"
		fmt.Println(msg)
		return
	}
	logicRun[nmRoutine] = logic

	return
}

// ResetLogic - ResetLogic
func ResetLogic(nmRoutine string) {
	delete(logicRun, nmRoutine)
}
