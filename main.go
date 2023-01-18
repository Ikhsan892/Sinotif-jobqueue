package main

import (
	"github.com/adjust/rmq/v5"
	"log"
	"os"
	"os/signal"
	"sinotif/configs"
	"sinotif/internal/consumer"
	"sinotif/pkg/handlers"
	"sinotif/pkg/utils"
	"syscall"
	"time"
)

const (
	prefetchLimit = 1000
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 5

	reportBatchSize = 10000
	consumeDuration = time.Millisecond
	shouldLog       = true
)

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Print("consume error: ", err)
		case *rmq.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}

func main() {
	cfg := configs.NewConfig()
	db := consumer.ConnectionDB(cfg.DB)

	errChan := make(chan error)

	go logErrors(errChan)
	conn, errConnection := rmq.OpenConnectionWithRedisClient("sinotif consumer job", consumer.ConnectRedis(), errChan)
	if errConnection != nil {
		utils.Error(utils.DATABASE, errConnection)
		os.Exit(1)
	}

	// register new queue
	queue, errQueue := conn.OpenQueue("test")
	smallTalk, errSmallTalk := conn.OpenQueue("report_small_talk")
	if errSmallTalk != nil {
		panic(errSmallTalk)
	}
	if errQueue != nil {
		panic(errQueue)
	}

	// start consuming
	if errConsuming := queue.StartConsuming(prefetchLimit, pollDuration); errConsuming != nil {
		panic(errConsuming)
	}
	if errConsumingSmallTalk := smallTalk.StartConsuming(prefetchLimit, pollDuration); errConsumingSmallTalk != nil {
		panic(errConsumingSmallTalk)
	}

	// assign handlers
	go handlers.HandlerTest(queue)
	go handlers.HandlerReportSmallTalk(smallTalk, db, cfg)

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-sigQuit

	utils.Info(utils.DATABASE, "Closing Connection")
	<-conn.StopAllConsuming()
}
