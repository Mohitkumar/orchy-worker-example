package main

import (
	"fmt"
	"time"

	"github.com/mohitkumar/orchy/worker"
)

func main() {
	config := &worker.WorkerConfiguration{
		ServerUrl:                []string{"localhost:8099", "localhost:8098"},
		MaxRetryBeforeResultPush: 1,
		RetryIntervalSecond:      1,
	}

	addParamActionFn := func(data map[string]any) (map[string]any, error) {
		data["newKey"] = 22
		return data, nil
	}
	logActionFn := func(data map[string]any) (map[string]any, error) {
		fmt.Println(data)
		return data, nil
	}
	tp := worker.NewWorkerConfigurer(*config)

	addParamWorker := worker.NewDefaultWorker(addParamActionFn).WithRetryCount(1).WithTimeoutSeconds(20)
	logWorker := worker.NewDefaultWorker(logActionFn)

	err := tp.RegisterWorker(addParamWorker, "add-data-worker", 1*time.Second, 2, 1)
	fmt.Print(err)
	tp.RegisterWorker(logWorker, "print-worker", 1*time.Second, 2, 1)
	tp.Start()
}
