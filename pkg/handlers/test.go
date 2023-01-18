package handlers

import (
	"fmt"
	"github.com/adjust/rmq/v5"
	"sinotif/pkg/functions"
)

func HandlerTest(queue rmq.Queue) {
	// maximum 5 worker or more
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("worker test consumer %d", i)
		if _, err := queue.AddConsumer(name, functions.NewTest(i)); err != nil {
			panic(err)
		}
	}
}
