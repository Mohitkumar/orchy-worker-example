package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mohitkumar/orchy/worker"
)

func main() {
	var wg sync.WaitGroup
	config := &worker.WorkerConfiguration{
		ServerUrl: "localhost:8099",
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
	tp := worker.NewWorkerConfigurer(*config, &wg)

	addParamWorker := worker.NewDefaultWorker(addParamActionFn).WithRetryCount(2).WithTimeoutSeconds(20)
	logWorker := worker.NewDefaultWorker(logActionFn).WithRetryCount(2).WithTimeoutSeconds(20)
	enhanceDataWorker := worker.NewDefaultWorker(enhanceDataFn).WithRetryCount(2).WithTimeoutSeconds(20)

	tp.RegisterWorker(addParamWorker, "add-data-worker", 100*time.Millisecond, 100, 1)
	tp.RegisterWorker(logWorker, "print-worker", 100*time.Millisecond, 100, 1)
	tp.RegisterWorker(enhanceDataWorker, "enhanceData", 100*time.Millisecond, 100, 1)
	tp.Start()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	fmt.Println("stopping")
	tp.Stop()
	wg.Wait()
}
