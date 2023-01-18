package handlers

import (
	"fmt"
	"github.com/adjust/rmq/v5"
	"gorm.io/gorm"
	"sinotif/configs"
	"sinotif/pkg/functions"
)

func HandlerReportSmallTalk(queue rmq.Queue, db *gorm.DB, cfg *configs.Config) {
	// maximum 5 worker or more
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("worker report small talk consumer %d", i)
		if _, err := queue.AddConsumer(name, functions.NewReportSmalltalk(db, i, cfg)); err != nil {
			panic(err)
		}
	}
}
