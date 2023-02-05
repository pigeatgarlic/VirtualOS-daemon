package fabric

import (
	"sync"

	childprocess "github.com/pigeatgarlic/VirtualOS-daemon/child-process"
	ws "github.com/pigeatgarlic/VirtualOS-daemon/fabric/wsocket"
)

type Command struct {
	Process string `json:"process"`
	Command string `json:"command"`
	Data    string `json:"data"`
}

type Fabric struct {
	mutex *sync.Mutex
	tenants map[childprocess.ProcessID]ws.IwebsocketTenant

	childprocess *childprocess.ChildProcesses
	wsserver *ws.WebsocketFabric
}



func NewFaric(childprocesssys *childprocess.ChildProcesses) (ret *Fabric,err error){
	ret = &Fabric {
		childprocess: childprocesssys,
		tenants: make(map[childprocess.ProcessID]ws.IwebsocketTenant),
		mutex: &sync.Mutex{},
	}


	ret.wsserver,err = ws.NewServer("localhost:3000","/fabric",func(conn ws.IwebsocketTenant, secret string) {
		ret.mutex.Lock()
		defer ret.mutex.Unlock()

		ProcessID := ret.childprocess.FindIDfromSecret(secret)
		ret.tenants[ProcessID] = conn
	})

	if err !=nil {
		return nil,err	
	}
	return
}


