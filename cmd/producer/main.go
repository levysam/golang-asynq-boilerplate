package main

import (
	"encoding/json"
	"fmt"
	"golang-queue-boilerplate/workers/sample"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const redisAddr = "localhost:6381"

const (
	TypeImageResize = "sample:queue"
)

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	for i := 0; i < 1000; i++ {
		task, err := NewImageResizeTask("http://image.com/" + fmt.Sprint(i))
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		info, err := client.Enqueue(task)
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(sample.SamplePayload{Data: src})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}
