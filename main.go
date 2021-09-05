package main

import (
	"flag"
	"fmt"
	"log"
	"marcode.io/echo-server/server"
	"os"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Run()
	}()

	fmt.Printf("TCP Echo Server is running on port %s \n", *port)

	wg.Wait()
}
