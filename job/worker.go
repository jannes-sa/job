package job

import (
	"fmt"
)

func sendinput(
	mappingTasks map[string]map[int]interface{},
	nmRoutine string,
	input chan<- interface{},
) {
	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			fmt.Printf("write: error writing %d on channel: %v\n", input, err)
			return
		}
	}()

	for k, v := range mappingTasks[nmRoutine] {
		in := correlated{
			key:   k,
			input: v,
		}
		input <- in
		delete(mappingTasks[nmRoutine], k)
	}

}

func getOutput(countTasks int, nmRoutine string, output chan correlated) {
	var out OutputData
	out.Tasks = nmRoutine

	for i := 0; i < countTasks; i++ {
		out.TotalTasks++

		o, ok := <-output
		if ok {
			out.Result = append(out.Result, o.output)
			if o.err != nil {
				out.Err = append(out.Err, WrapperOutputError{
					Err:        o.err,
					InputError: o.input,
				})
				out.TotalTasksFail++
			} else {
				out.TotalTasksDone++
			}
		}
	}
	out.TotalTasksPending = countTasks - len(out.Result)

	if !logicRun[nmRoutine].Done(&out) {
		mappingStatusTasks[nmRoutine] = restart
		return
	}
	mappingStatusTasks[nmRoutine] = done
}

func tearDown(
	input chan<- interface{},
	output chan<- correlated,
) {
	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			fmt.Printf("write: error writing %d on channel: %v\n", input, err)
			return
		}
	}()

	close(input)
	close(output)
}

// ChanInputData - Channel Receiver Data
type ChanInputData struct {
	Data interface{}
}
type OutputData struct {
	Tasks             string               // Task Name
	TotalTasks        int                  // Total Task Has Been Received By Output
	TotalTasksDone    int                  // Total Task Has Been Received By Output And Done According logic
	TotalTasksFail    int                  // Total Task Has Been Received By Output And Fail According logic
	TotalTasksPending int                  // Total Task Has Not Received By Output
	Result            []interface{}        // All Wrapping Result
	Err               []WrapperOutputError // All Error Contain inside of this val
}

type WrapperOutputError struct {
	Err        error
	InputError interface{}
}

type correlated struct {
	key    int
	input  interface{}
	output interface{}
	err    error
}

func worker(
	input chan interface{},
	output chan correlated,
	nmRoutine string,
) {

	for data := range input {
		d := data.(correlated)
		d.output, d.err = logicRun[nmRoutine].Run(ChanInputData{
			Data: d.input,
		})
		output <- d
	}
}
