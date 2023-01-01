package main

import (
	"fmt"
	"time"

	"github.com/mohitkumar/orchy/worker"
)

func main() {
	config := &worker.WorkerConfiguration{
		ServerUrl:                "localhost:8099",
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
	enhanceDataFn := func(data map[string]any) (map[string]any, error) {
		data["key"] = fmt.Sprintf("prefix_%v", data["key"])
		return data, nil
	}
	tp := worker.NewWorkerConfigurer(*config)

	addParamWorker := worker.NewDefaultWorker(addParamActionFn).WithRetryCount(2).WithTimeoutSeconds(20)
	logWorker := worker.NewDefaultWorker(logActionFn).WithRetryCount(2).WithTimeoutSeconds(20)
	enhanceDataWorker := worker.NewDefaultWorker(enhanceDataFn).WithRetryCount(2).WithTimeoutSeconds(20)

	tp.RegisterWorker(addParamWorker, "add-data-worker", 1*time.Millisecond, 100, 4)
	tp.RegisterWorker(logWorker, "print-worker", 1*time.Millisecond, 100, 4)
	tp.RegisterWorker(enhanceDataWorker, "enhanceData", 1*time.Millisecond, 100, 4)
	tp.Start()
}
