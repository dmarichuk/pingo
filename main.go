package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pingo/parser"
	"syscall"
)

var (
	pathToConfig *string
)

func init() {
	pathToConfig = flag.String("config", "./pingo.yaml", "path to config")
}

func main() {
	flag.Parse()
	parser.YamlConfig.ReadFromFile(*pathToConfig)
	jobs := parser.YamlConfig.Parse()

	log.Println("JOBS: ", jobs)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	for i := 0; i < len(jobs); i++ {
		log.Println("Current job", jobs[i].Name)
		go jobs[i].Run()
	}

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
