package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pingo/dashboard"
	db "pingo/database"
	"pingo/parser"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

var (
	pathToConfig  *string
	dashboardPort *string
)

func init() {
	pathToConfig = flag.String("config", "/pingo/config/pingo.yaml", "path to config")
	dashboardPort = flag.String("port", "9080", "port for dashboard")
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
		dbErr := jobs[i].DumpToDB(db.DB)
		if dbErr != nil {
			log.Fatal(dbErr)
		}
		go jobs[i].Run()
	}

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
		os.Exit(1)
	}()

	if err := http.ListenAndServe(":"+*dashboardPort, dashboard.Mux); err != nil {
		log.Fatal(err)
	}

}
