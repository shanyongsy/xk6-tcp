package tcp

import (
	"bufio"
	"fmt"
	"net"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/tcp", new(TCP))
}

type CallBack func(data string)

type TCP struct {
	conn     net.Conn
	connStr  string
	lastErr  error
	onRevMsg CallBack
}

func (tcp *TCP) Connect(addr string, onRevMsg CallBack) error {
	tcp.connStr = addr
	tcp.onRevMsg = onRevMsg
	tcp.conn, tcp.lastErr = net.Dial("tcp", tcp.connStr)
	if tcp.lastErr != nil {
		tcp.conn = nil
		return tcp.lastErr
	} else {
		go tcp.readConn()
	}

	return nil
}

func (tcp *TCP) Write(data []byte) error {
	_, err := tcp.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (tcp *TCP) WriteLn(data []byte) error {
	_, err := tcp.conn.Write(append(data, []byte("\n")...))
	return err
}

func (tcp *TCP) readConn() {
	for {
		scanner := bufio.NewScanner(tcp.conn)

		for {
			ok := scanner.Scan()
			text := scanner.Text()

			if !ok {
				fmt.Println("Reached EOF on server connection.")
				break
			} else {
				tcp.onRevMsg(text)
			}
		}
	}
}
