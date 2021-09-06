package server

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"testing"
)

func TestTCPEchoServer(t *testing.T) {
	tt := []struct {
		name          string
		clientMessage string
		serverReply   string
	}{
		{"Client says, Hello World, and Server replies with same message", "Hello World\n", "Hello World\n"},
		{"Client says, Goodbye World, and Server replies with same message", "Goodbye World\n", "Goodbye World\n"},
	}

	server := NewTCPServer(":1235")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cert, err := os.ReadFile("/home/marcode/tls" + "/cert.pem")
			if err != nil {
				t.Fatal(err)
			}
			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM(cert); !ok {
				t.Fatalf("unable to parse cert")
			}
			config := &tls.Config{RootCAs: certPool}

			conn, err := tls.Dial("tcp", "localhost:1235", config)
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					t.Fatalf("error closing the connection: %v", err)
				}
			}(conn)

			if _, err := conn.Write([]byte(tc.clientMessage)); err != nil {
				t.Error("client could not write message to TCP server:", err)
			}

			reader := bufio.NewReader(conn)
			response, err := reader.ReadString(byte('\n'))
			if err != nil {
				t.Fatalf("could not read from connection")
			}

			if response != tc.serverReply {
				t.Errorf("expected the following reply from server %q, but got %q", tc.serverReply, response)
			}
		})
	}

	server.Stop()
}
