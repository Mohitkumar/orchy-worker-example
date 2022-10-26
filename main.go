package main

import (
	"fmt"
	"time"

	"github.com/mohitkumar/orchy/worker"
)

func main() {
	config := &worker.WorkerConfiguration{
		ServerUrl:                "localhost:8099",
		PollInterval:             1,
		MaxRetryBeforeResultPush: 1,
		RetryIntervalSecond:      1,
	}

	addDataWorkerFn := func(data map[string]any) (map[string]any, error) {
		data["newKey"] = 22
		return data, nil
	}
	printWorkerFn := func(data map[string]any) (map[string]any, error) {
		fmt.Println(data)
		return data, nil
	}
	tp := worker.NewWorkerConfigurer(*config)

	addWorker := worker.NewDefaultWorker(addDataWorkerFn).WithRetryCount(1).WithTimeoutSeconds(20)
	printWorker := worker.NewDefaultWorker(printWorkerFn)

	tp.RegisterWorker(addWorker, "add-data-worker", 1*time.Second, 2, 1)
	tp.RegisterWorker(printWorker, "print-worker", 1*time.Second, 2, 1)
	tp.Start()
}
