package connectionPool

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionManager struct {
	TCPListener   net.Listener
	Connections   map[net.Conn]struct{}
	mu            sync.Mutex
	addChannel    chan net.Conn
	removeChannel chan net.Conn
	MiddleWare    []func(net.Conn) error
}

func NewConnectionManager(size int) *ConnectionManager {
	return &ConnectionManager{
		Connections:   make(map[net.Conn]struct{}, size),
		addChannel:    make(chan net.Conn, size),
		removeChannel: make(chan net.Conn, size),
	}

}

func (c *ConnectionManager) AssignTCPListener(listener net.Listener) {
	c.TCPListener = listener
}

func (c *ConnectionManager) AddChannel(conn net.Conn) {
	c.addChannel <- conn
}
func (c *ConnectionManager) RemoveChannel(conn net.Conn) {
	c.removeChannel <- conn
}

func (c *ConnectionManager) Run() {
	for {
		select {
		case conn := <-c.addChannel:
			c.mu.Lock()
			fmt.Printf("Added IP Address Channel: %v\n", conn.RemoteAddr().String())
			c.Connections[conn] = struct{}{}
			c.mu.Unlock()
		case conn := <-c.removeChannel:
			c.mu.Lock()
			fmt.Printf("Removed  IP Address Channel: %v\n", conn.RemoteAddr().String())
			delete(c.Connections, conn)
			c.mu.Unlock()
		}

	}
}

func (c *ConnectionManager) ProcessConnections() {
	for {
		conn, err := c.TCPListener.Accept()
		if err != nil {
			fmt.Print("TCP accept failed.\n")
			return
		}

		go c.AddChannel(conn)

		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			go c.RemoveChannel(conn)
			break
		}

		c.distributePacketConn(conn, buffer)

	}
}

func (p *ConnectionManager) AddMiddleware(mw func(net.Conn) error) *ConnectionManager {
	p.MiddleWare = append(p.MiddleWare, mw)
	return p
}

func (c *ConnectionManager) distributePacketConn(distConn net.Conn, packet []byte) {
	for conn := range c.Connections {
		if conn.RemoteAddr().String() == distConn.RemoteAddr().String() {
			continue
		}

		conn.Write(packet)
	}
}
