package main

import (
	"flag"
	"log"
	"marcode.io/echo-server/server"
	"os"
	"os/signal"
)

func main() {

	protocol := flag.String("protocol", "tcp", "The IP Protocol for the server. Default is TCP")
	port := flag.String("port", os.Getenv("TCP_PORT"), "Port the server listens to. Default 1234")
	flag.Parse()

	address := ":" + *port

	srv, err := server.New(*protocol, address)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Received terminate, server is shutting down, signal: ", sig)

	srv.Stop()

	log.Println("Server Stopped")
}
