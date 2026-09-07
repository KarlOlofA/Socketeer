package network

import (
	"net"

	"github.com/KarlOlofA/socketeer/internal/auth"
)

type Connection struct {
	User       auth.User
	Connection net.Conn
}
