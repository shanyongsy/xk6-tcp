package tcp

import (
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"
	"reflect"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/shanyongsy/xk6-tcp/pb"
)

func init() {
	modules.Register("k6/x/tcp", new(TCP))
}

// TCP is the k6 tcp extension.
type TCP struct{}

// Client is the TCP Client wrapper.
type Client struct {

	// conn
	conn net.Conn

	// address and port, for example, 127.0.0.1:8000.
	connStr string

	// The last error from this struct.
	lastErr error

	// The function pointer used to receive the message needs to be passed in JS.
	onJSCallBack JSCallBack

	// The uuid of client
	uuid string

	// The token of client
	token string
}

// XClient represents the Client constructor (i.e. `new tcp.Client()`) and
// returns a new TCP client object.
func (r *TCP) XClient(ctxPtr *context.Context) interface{} {
	rt := common.GetRuntime(*ctxPtr)
	return common.Bind(rt, &Client{}, ctxPtr)
}

// Create new tcp.
// A new TCP link must be used, otherwise it will cause the same conn to communicate.
//func (tcp *Client) Create() *Client {
//	return new(Client)
//}

// Send the received message to JS. Com through this function.
// cmd - msg id
// sus -resault
type JSCallBack func(cmd string, sus bool)

// To init all things.
func (client *Client) Connect(addr string, onJSCallBack JSCallBack) error {
	client.connStr = addr
	client.uuid = uuid.New().String()
	client.token = client.uuid
	client.onJSCallBack = onJSCallBack
	client.conn, client.lastErr = net.Dial("tcp", client.connStr)
	if client.lastErr != nil {
		client.conn = nil
		return client.lastErr
	} else {
		go client.loopReadConn()
	}

	return nil
}

// Send msg by this function.
// func (client *Client) WriteStr(data string) error {

// 	if client.conn == nil {
// 		return errors.New("call Write function, but conn is nil.")
// 	}

// 	_, err := client.conn.Write([]byte(data))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Send msg by this function.
// func (client *Client) WriteStrLn(data string) error {
// 	return client.WriteStr(fmt.Sprintln(data))
// }

// Get msg by this function.
func (client *Client) loopReadConn() {
	defer func() {
		log.Println("loopReadConn out.")
	}()

	for {
		lenBuf := make([]byte, 2)
		_, err := io.ReadFull(client.conn, lenBuf)
		if err != nil {
			if err == io.EOF {
				log.Println("client quit...")
			} else {
				log.Println("read length error: " + err.Error())
			}
			break
		}
		pkgLen := binary.BigEndian.Uint16(lenBuf[0:])
		buf := make([]byte, pkgLen)
		_, err = io.ReadFull(client.conn, buf)
		if err != nil {
			log.Println("read package error: " + err.Error())
			break
		}
		onRecveMsg(client, &buf)
	}
}

func (client *Client) sendInfoToJS(cmd *string, sus bool) {
	if client.onJSCallBack != nil {
		client.onJSCallBack(*cmd, sus)
	}
}

func (client *Client) sendMsgToGateway(cmd *string, body *[]byte) error {

	frameMsg := &pb.FrameMessage{Cmd: *cmd, Body: *body}
	msg := &pb.C2GMessage{GatewayMessage: frameMsg, BusinessMessage: nil}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	buffer := make([]byte, len(data)+2)

	binary.BigEndian.PutUint16(buffer[0:], uint16(len(data)))
	copy(buffer[2:], data)
	client.conn.Write(buffer)

	return nil
}

func onRecveMsg(client *Client, data *[]byte) {
	var cmd string
	var sus bool = false

	msg := &pb.G2CMessage{}
	err := proto.Unmarshal(*data, msg)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if msg.GatewayMessage != nil {
		cmd, sus = onDealMsg(client, msg.GatewayMessage)
	} else {
		return
	}

	client.sendInfoToJS(msgMapping(&cmd), sus)
}

// 返回 - param1:G2C cmdID, param2:sus
func onDealMsg(client *Client, msg *pb.FrameMessage) (string, bool) {

	value := reflect.ValueOf(client)
	f := value.MethodByName(msg.Cmd)
	if !f.IsValid() || f.Kind() != reflect.Func {
		log.Printf("No handler found for message %v.\n", msg.Cmd)
		return msg.Cmd, false
	}

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(&msg.Body)
	ret := f.Call(in)

	if ret == nil || len(ret) != 1 || ret[0].Kind() != reflect.Bool {
		log.Printf("The processing function of message %v is incorrectly written.\n", msg.Cmd)
		return msg.Cmd, false
	}

	return msg.Cmd, ret[0].Interface().(bool)
}

// 返回收发消息的对应关系
func msgMapping(recvMsg *string) *string {
	return recvMsg
}
