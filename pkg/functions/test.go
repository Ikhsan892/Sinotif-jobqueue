package functions

import (
	"fmt"
	"github.com/adjust/rmq/v5"
	"time"
)

type Test struct {
	workerIndex int
}

func NewTest(workerIndex int) *Test {
	return &Test{
		workerIndex: workerIndex,
	}
}

func (t *Test) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()

	time.Sleep(2 * time.Second)
	fmt.Println(fmt.Sprintf("execute worker test index %d with payload %s", t.workerIndex, payload))
	delivery.Reject()
}
