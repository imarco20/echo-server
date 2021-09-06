package server

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type TCPServer struct {
	address  string
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func NewTCPServer(addr string) *TCPServer {
	server := &TCPServer{
		address: addr,
		quit:    make(chan interface{}),
	}

	cert, err := tls.LoadX509KeyPair(os.Getenv("certs")+"/cert.pem", os.Getenv("certs")+"/key.pem")

	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", addr, &config)
	if err != nil {
		log.Fatal(err)
	}

	server.listener = listener
	server.wg.Add(1)

	go server.Run()

	log.Printf("TCP Server is running on Port %s", addr)

	return server
}

func (t *TCPServer) Run() {

	defer t.wg.Done()

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.quit:
				return
			default:
				log.Println("error accepting client connection", err)
			}

		} else {
			t.wg.Add(1)
			go func() {
				t.handleConnection(conn)
				t.wg.Done()
			}()
		}

	}
}

func (t *TCPServer) Stop() {
	close(t.quit)
	err := t.listener.Close()
	if err != nil {
		log.Println("error closing connection", err)
	}
	t.wg.Wait()
}

func (t *TCPServer) handleConnection(conn net.Conn) {

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	buf := make([]byte, 1024)

ReadLoop:
	for {
		select {
		case <-t.quit:
			return
		default:
			err := conn.SetDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				log.Println("error setting a read and write deadline to the connection", err)
			}
			n, err := conn.Read(buf)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue ReadLoop
				} else if err != io.EOF {
					log.Println("encountered an error while reading from connection", err)
					return
				}

			}

			if n == 0 {
				return
			}

			log.Printf("received from client: %s", string(buf[:n]))
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Println("error writing response to client", err)
				return
			}

			log.Printf("sent from server: %s", string(buf[:n]))
		}

	}
}
