package edgerouter

import (
	"errors"
	"fmt"
	"net"
)

type TcpClient struct {
	ctl   Controller
	Port  int
	conns map[string]*net.TCPConn
}

func (c *TcpClient) Init() {
	c.conns = make(map[string]*net.TCPConn)
}

func (c *TcpClient) PacketReceived(bts []byte, conn *net.TCPConn) int {
	return c.ctl.OnReceived(bts, conn)
}

func (c *TcpClient) SetController(ctrl Controller) {
	c.ctl = ctrl
}

func (c *TcpClient) Connect(to string) error {
	if addr, err := net.ResolveTCPAddr("tcp", to); err == nil {
		var conn *net.TCPConn
		var ok bool
		if conn, ok = c.conns[addr.String()]; !ok {
			var localPort *net.TCPAddr
			if c.Port != 0 {
				localPort, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.Port))
			}
			if conn, err = net.DialTCP("tcp", localPort, addr); err == nil {
				go handleTcpConn(conn, c)
				c.conns[addr.String()] = conn
			}
		}
		return err
	} else {
		return err
	}
}

func (c *TcpClient) Send(msg *BytesMessage) (err error) {
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", msg.To); err == nil {
		if conn, ok := c.conns[addr.String()]; ok {
			_, err = conn.Write(msg.Message)
		} else {
			err = errors.New("no such connection")
		}
	}

	return err
}

func (c *TcpClient) String() string {
	return fmt.Sprintf("tcp client(%p)", c)
}