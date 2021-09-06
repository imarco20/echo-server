# TCP Echo Server

An echo service is a useful debugging and measurement tool. It's specified in [RFC 862](https://datatracker.ietf.org/doc/html/rfc862). A TCP echo server listens for TCP
connections on a specific TCP port. Once a connection is established, the echo server sends back any data it receives.
This continues until the calling user terminates the connection.

## How to run the Echo Server

1. Using Terminal

- Navigate to the directory of the application on your computer, and run the following command

```
go run . -port=port_of_your_choice
```

- Then open another terminal and connect to the server using the following command

```
openssl s_client -connect localhost:port_server_listens_to
```

2. Using Docker

- You can build a docker image for the application using the Dockerfile available in the project's root directory.
- Build the image by running the following command from your terminal:

```
docker build -t tcp-echo-server:latest .
```

- Then to create and run a docker container out of the image, enter the command:

```
docker run -d --rm -it -e certs="build/tls" -e TCP_PORT=port_of_your_choice -p port_of_your_choice:port_of_your_choice tcp-echo-server:latest
```

- Then open another terminal and connect to the server using the following command

```
openssl s_client -connect localhost:port_server_listens_to
```