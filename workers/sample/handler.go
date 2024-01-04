package sample

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-queue-boilerplate/pkg/registry"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type SamplePayload struct {
	Data string
}

type SampleWorker struct {
	reg *registry.Registry
}

func (processor *SampleWorker) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p SamplePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("payload data string = %s", p.Data)
	// Code
	time.Sleep(2 * time.Second)
	return nil
}

func NewImageProcessor(reg *registry.Registry) *SampleWorker {
	return &SampleWorker{
		reg: reg,
	}
}
