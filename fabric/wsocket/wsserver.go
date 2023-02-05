package ws

import (
	"fmt"

	"golang.org/x/net/websocket"
)


type IwebsocketTenant interface {
	Read() (string,error)
	Write(string)
	Exited() bool
}


type websocketTenant struct {
	conn *websocket.Conn

	exited bool
}


func (ws *websocketTenant)	Read() (string,error) {
	if ws.exited {
		return "",fmt.Errorf("tenant exited")
	}

	bytes := make([]byte,1024)
	n,err := ws.conn.Read(bytes)
	if err != nil {
		ws.exited = true
		return "",err
	}
	return string(bytes[:n]),nil
}

func (ws *websocketTenant)	Write(data string) {
	if ws.exited {
		return
	}

	_,err := ws.conn.Write([]byte(data))	
	if err != nil {
		ws.exited = true
	}
}
func (ws *websocketTenant)	Exited()bool {
	return ws.exited
}

type WebsocketFabric struct {
	server websocket.Server

	ClientCallback func(conn IwebsocketTenant,secret string)
}


func NewServer(Adress string, Path string, callback func(conn IwebsocketTenant, secret string)) (*WebsocketFabric, error){
	Config,err := websocket.NewConfig(fmt.Sprintf("ws://%s%s", Adress, Path), "http://localhost")
	if err != nil {
		return nil,err;
	}

	ws := &WebsocketFabric{
		ClientCallback: callback,
		server: websocket.Server{
			Config: *Config,
		},
	}
	ws.server.Handler = ws.EchoServer
	return ws,nil
}

func (ws *WebsocketFabric)EchoServer(con *websocket.Conn) {
	auth := con.Request().Header.Get("Authorization")
	tenant := &websocketTenant{
		conn: con,
		exited: false,
	}
	ws.ClientCallback(tenant,auth)	
}



// func (procs *ChildProcesses) GetIncomingMessage(ID ProcessID) string {
// 	procs.mutex.Lock()
// 	defer procs.mutex.Unlock()

// 	// TODO
// 	return ""
// }

// func (procs *ChildProcesses) SendMessage(ID ProcessID, val string){
// 	procs.mutex.Lock()
// 	defer procs.mutex.Unlock()

// 	// TODO
// }
