package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	db "pingo/database"
	"pingo/parser"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	defer db.DB.Close()

	for i := 0; i < len(jobs); i++ {
		go jobs[i].Run()
	}

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
