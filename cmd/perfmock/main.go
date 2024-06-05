package main

import (
	"fmt"
	"time"

	"github.com/oke11o/perf-mock/internal/http"
	"github.com/oke11o/perf-mock/internal/logger"
)

func main() {
	addr := "localhost:8091"
	log := logger.New()
	server := http.NewServer(addr, log, time.Now().UnixNano())
	server.ServeAsync()

	err := <-server.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("FINISH")
}
