package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net"

	pool "github.com/KarlOlofA/socketeer/internal/network"
	types "github.com/KarlOlofA/socketeer/internal/types/network"
)

func AuthMiddleware(key string) pool.Middleware {
	return func(conn net.Conn) (net.Conn, error) {
		p := types.Packet{}
		var buffer bytes.Buffer
		_, err := io.Copy(&buffer, conn)
		if err != nil {
			return nil, err
		}

		err = p.FromByteSlice(buffer.Bytes())
		if err != nil {
			return nil, err
		}

		if p.Key != key {
			return nil, fmt.Errorf("invalid key: expected %s, got %s", key, p.Key)
		}
		return conn, nil
	}
}
