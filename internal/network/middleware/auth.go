package middleware

import (
	"fmt"
	"net"

	pool "github.com/KarlOlofA/socketeer/internal/network"
	types "github.com/KarlOlofA/socketeer/internal/types/network"
)

func AuthMiddleware(key string) pool.Middleware {
	return func(conn net.Conn) (net.Conn, error) {
		p := types.Packet{}
		var buffer []byte = make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			return nil, err
		}

		err = p.FromByteSlice(buffer)
		if err != nil {
			return nil, err
		}

		if p.Key != key {
			return nil, fmt.Errorf("invalid key: expected %s, got %s", key, p.Key)
		}
		return conn, nil
	}
}
