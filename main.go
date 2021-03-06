package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"statusServer/statuscheck"
)

func runServer() bool {
	useHttp := flag.Bool("useHttp", false, "Run useHttp server with statuses")
	flag.Parse()
	return *useHttp
}

func getFilename() string {
	if flag.NArg() != 1 {
		log.Fatal("Please provide config for services to check")
	}
	return flag.Arg(0)
}

func cmdHandler(filename string) {
	config := statuscheck.ReadConfig(filename)
	status := statuscheck.PingServices(config)
	fmt.Println("Service status check")
	for i := 0; i < len(status); i++ {
		stat := status[i]
		if stat.IsRunning {
			fmt.Printf("\t%s...OK\n", stat.ServiceName)
		} else {
			fmt.Printf("\t%s...error\n\t\t%s\n", stat.ServiceName, stat.Error)
		}
	}
}

func httpHandler(w http.ResponseWriter, _ *http.Request) {
	filename := getFilename()
	config := statuscheck.ReadConfig(filename)
	status := statuscheck.PingServices(config)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(status)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = w.Write(js)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	if runServer() {
		http.HandleFunc("/", httpHandler)
		log.Fatal(http.ListenAndServe(":5555", nil))
	} else {
		cmdHandler(getFilename())
	}
}
