package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"post-it/pkg/client"
	"post-it/pkg/controller"
)

// Register log formatting
func init() {
	formatter := log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
	}
	log.SetFormatter(&formatter)
}

func main() {
	method := flag.String("method", "GET", "Method")
	url := flag.String("url", "", "Url")
	input := flag.String("input", "input.csv", "Input File")
	output := flag.String("output", "output.csv", "Output File")
	batchSize := flag.Int("batchSize", 100, "Batch Size")
	threadPool := flag.Int("threadPool", 10, "Thread Pool")

	flag.Parse()

	if url == nil || *url == "" {
		flag.PrintDefaults()
		return
	}

	run(*method, *url, *input, *output, *batchSize, *threadPool)
}

func run(method, url, input, output string, batchSize, threadPool int){
	// API CLIENT
	clt, err := client.NewClient(client.Config {
		Timeout: 3000,
		InsecureSkipVerify: true,
		MaxConnsPerHost: 100,
		MaxIdleConns: 10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout: 5000,
	})
	if err != nil {
		log.Fatal(err)
	}

	// CONTROLLER
	ctrl := controller.Controller{Method: method, Url: url, Client: clt, BatchSize: batchSize, ThreadPool: threadPool}
	err = ctrl.Run(input, output)
	if err != nil {
		log.Error(err)
	}
}
