package main

import (
	"fmt"
	"log"
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
    
    for i:=0; i < len(jobs); i++ {
        log.Println("Current job", jobs[i].Name)
        go jobs[i].RunJob()
	}


	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
