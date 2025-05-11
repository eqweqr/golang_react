package server

import (
	"database/sql"
	"fmt"
	"net"
)

type Server struct {
	Secret     string
	Ip_address string
	Port       int
	DB         *sql.DB
}

func NewServer(secret string, ip string, port int, db *sql.DB) (*Server, error) {
	if ip := net.ParseIP(ip); ip == nil {
		return nil, fmt.Errorf("%s not valid ip address", ip)
	}
	return &Server{
		Secret:     secret,
		Ip_address: ip,
		Port:       port,
		DB:         db,
	}, nil

}
