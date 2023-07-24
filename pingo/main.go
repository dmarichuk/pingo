package main

import (
	"fmt"
	"os"
	"os/signal"
	"pingo/parser"
	"syscall"
)

func main() {
	jobs := parser.YamlConfig.Parse()
	fmt.Println(jobs)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	for _, job := range jobs {
		go job.RunJob()
	}

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
