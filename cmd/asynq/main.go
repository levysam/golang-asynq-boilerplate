package main

import (
	"golang-queue-boilerplate/pkg/config"
	"golang-queue-boilerplate/pkg/logger"
	"golang-queue-boilerplate/pkg/registry"
	"golang-queue-boilerplate/workers/sample"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

var (
	reg *registry.Registry
)

const (
	SampĺeQueue = "sample:queue"
)

func main() {
	conf := reg.Inject("config").(*config.Config)
	redisHost := conf.ReadConfig("REDIS_HOST")
	redisPort := conf.ReadConfig("REDIS_PORT")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort},
		asynq.Config{
			Concurrency:     conf.ReadNumberConfig("CONCURRENCY"),
			ShutdownTimeout: 5 * time.Second,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(SampĺeQueue, sample.NewImageProcessor(reg))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func init() {
	conf := config.NewConfig()

	l := logger.NewLogger(
		conf.ReadConfig("ENVIRONMENT"),
		conf.ReadConfig("APP"),
		conf.ReadConfig("VERSION"),
	)

	l.Boot()

	reg = registry.NewRegistry()
	reg.Provide("config", conf)
	reg.Provide("logger", l)
}
