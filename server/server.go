package server

import (
	"errors"
	"strings"
)

type Server interface {
	Run()
	Stop()
}

func New(protocol, address string) (Server, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return NewTCPServer(address), nil
	default:
		return nil, errors.New("unsupported protocol")
	}
}
