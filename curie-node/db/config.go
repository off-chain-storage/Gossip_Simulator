package db

import (
	"context"
	"net"
)

type Config struct {
	DbAddr       string
	PoolSize     uint
	MaxIdleConns uint
	PoolFIFO     bool
	Dialer       func(ctx context.Context, network, address string) (net.Conn, error)
}
