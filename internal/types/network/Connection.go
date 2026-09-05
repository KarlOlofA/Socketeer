package network

import (
	"net"

	"socketeer.github.com/internal/auth"
)

type Connection struct {
	User       auth.User
	Connection net.Conn
}
