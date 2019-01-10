package scheduler

import (
	"time"
	"util/channel/lib/scheduler/example/logic1"
)

func scheduler() {
	logic1.RunScheduler()
	// logic2.RunScheduler()

	time.Sleep(10 * time.Minute)
}
