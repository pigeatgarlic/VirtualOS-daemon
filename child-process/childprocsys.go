package childprocess

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/pigeatgarlic/VirtualOS-daemon/log"
	"golang.org/x/net/websocket"
)

type ProcessID int


type ChildProcess struct {
	exited bool
	cmd *exec.Cmd
}

type ChildProcesses struct {
	count int
	mutex sync.Mutex

	procs map[ProcessID]*ChildProcess
	server *websocket.Server
}

func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

// Origin parses the Origin header in req.
// If the Origin header is not set, it returns nil and nil.
func Origin(config *websocket.Config, req *http.Request) (*url.URL, error) {
	var origin string
	switch config.Version {
	case websocket.ProtocolVersionHybi13:
		origin = req.Header.Get("Origin")
	}
	if origin == "" {
		return nil, nil
	}
	return url.ParseRequestURI(origin)
}
func checkOrigin(config *websocket.Config, req *http.Request) (err error) {
	config.Origin, err = Origin(config, req)
	
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}

func NewChildProcessSystem(Adress string, Path string) (ret *ChildProcesses,err error) {
	ret = &ChildProcesses{
		procs: make(map[ProcessID]*ChildProcess),
		mutex: sync.Mutex{},
		count: 0,
	}

	Config,err := websocket.NewConfig(fmt.Sprintf("ws://%s%s", Adress, Path), "http://localhost")
	if err != nil {
		return nil,err;
	}

	ret.server = &websocket.Server{
		Handler: EchoServer, 
		Handshake: checkOrigin,
		Config: *Config,
	}
	return
}

func (procs *ChildProcesses) CloseAll() {
	procs.mutex.Lock()
	defer procs.mutex.Unlock()

	for ID,proc := range procs.procs {
		log.PushLog("force terminate process name %s, process id %d \n", proc.cmd.Args[0], int(ID))
		proc.cmd.Process.Kill()
	}
}

func (procs *ChildProcesses) CloseID(ID ProcessID) {
	procs.mutex.Lock()
	defer procs.mutex.Unlock()

	proc := procs.procs[ID]
	if proc == nil {
		return
	}

	log.PushLog("force terminate process name %s, process id %d \n", proc.cmd.Args[0], int(ID))
	proc.cmd.Process.Kill()
}

func (procs *ChildProcesses) GetName(ID ProcessID) string {
	procs.mutex.Lock()
	defer procs.mutex.Unlock()

	filename := strings.Split(procs.procs[ID].cmd.Args[0],"/")
	return strings.Split(filename[len(filename) - 1],".")[0]
}

func (procs *ChildProcesses) GetIncomingMessage(ID ProcessID) string {
	procs.mutex.Lock()
	defer procs.mutex.Unlock()

	// TODO
	return ""
}

func (procs *ChildProcesses) SendMessage(ID ProcessID, val string){
	procs.mutex.Lock()
	defer procs.mutex.Unlock()

	// TODO
}

func (procs *ChildProcesses) WaitID(ID ProcessID) {
	procs.mutex.Lock();
	proc := procs.procs[ID];
	procs.mutex.Unlock()

	if proc == nil {
		return
	}

	for {
		if !proc.exited {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
