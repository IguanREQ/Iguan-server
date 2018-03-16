package main

import (
	"iguan/listener"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"iguan/logs"
)

func Wait() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	logs.Info("OS Signal received: %v", <-signals)
}

func main() {
	go listener.RunHTTP(http.Server{
		Addr: "127.0.0.1:8080",
	})
	go listener.RunTCP("127.0.0.1:8081")

	// TODO: implement graceful shutdown on interrupts
	Wait()
}
