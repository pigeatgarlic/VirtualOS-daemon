package childprocess

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/pigeatgarlic/VirtualOS-daemon/log"
	"github.com/pigeatgarlic/VirtualOS-daemon/utils"
)

type ProcessID int
const (
	InvalidID = -1
)

type Process struct{
	exited bool
	cmd *exec.Cmd

	name string
	path string

	secret string
}


type ChildProcesses struct {
	count int
	mutex sync.Mutex
	procs map[ProcessID]*Process
}

func NewChildProcessSystem() (ret *ChildProcesses) {
	return &ChildProcesses{
		procs: make(map[ProcessID]*Process),
		mutex: sync.Mutex{},
		count: 0,
	}
}

func (procs *ChildProcesses) CloseAll() {
	for ID,_ := range procs.procs {
		procs.CloseID(ID)
	}
}


func (procs *ChildProcesses) NewChildProcess(dir *string, name string,args ...string) (ProcessID,error) {
	procs.mutex.Lock()
	defer func() {
		procs.count++
		procs.mutex.Unlock()
	}()


	
	path,err := utils.FindProcessPath(dir,name)
	if err != nil {
		return 0,err
	}
	cmd := exec.Command(path,args...)
	randstr := utils.CreateRandomString(20)
	cmd.Env = append(cmd.Env, fmt.Sprintf("DAEMON_SECRET=%s",randstr))

	id := ProcessID(procs.count)
	log.PushLog("process %s, process id %d booting up\n", name, int(id))
	procs.procs[id] = &Process{
		exited : false,
		cmd: cmd,
		name: name,
		path: path,
		secret: randstr,
	}

	go procs.handleProcess(id)
	return ProcessID(procs.count),nil
}




func (procs *ChildProcesses) FindIDfromSecret(secret string) ProcessID {
	procs.mutex.Lock()
	defer procs.mutex.Unlock()

	for pi, p := range procs.procs {
		if p.secret == secret {
			return pi
		}
	}
	return -1
}


func (procs *ChildProcesses) CloseID(ID ProcessID) {
	if !procs.findID(ID){
		return
	} 



	procs.mutex.Lock()
	proc := procs.procs[ID]
	procs.mutex.Unlock()
	if (proc.exited){
		return
	} else {
		log.PushLog("force terminate process name %s, process id %d \n", proc.name, int(ID))
		proc.cmd.Process.Kill()
	}
}


func (procs *ChildProcesses) WaitID(ID ProcessID) {
	if !procs.findID(ID) {
		return
	}

	procs.mutex.Lock();
	proc := procs.procs[ID];
	procs.mutex.Unlock()

	for {
		if !proc.exited {
			time.Sleep(100 * time.Millisecond)
		}
	}
}





