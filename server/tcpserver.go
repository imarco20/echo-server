package server

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type TCPServer struct {
	address  string
	listener net.Listener
}

func NewTCPServer(addr string) *TCPServer {
	server := &TCPServer{
		address: addr,
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

	return server
}

func (t *TCPServer) Run() {

	for {
		conn, err := t.listener.Accept()
		if err != nil {

			log.Println("error accepting client connection", err)
		} else {
			go func() {
				t.handleConnection(conn)
			}()
		}

	}
}

func (t *TCPServer) Stop() {
	err := t.listener.Close()
	if err != nil {
		log.Println("error closing connection", err)
	}
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
		conn.SetDeadline(time.Now().Add(5 * time.Second))
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
