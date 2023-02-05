package childprocess

import (
	"fmt"
	"io"
	"strings"

	"github.com/pigeatgarlic/VirtualOS-daemon/log"
)

func (procs *ChildProcesses) handleProcess(id ProcessID) {
	procs.mutex.Lock()
	proc := procs.procs[id]
	procs.mutex.Unlock()

	stdoutIn, _ := proc.cmd.StdoutPipe()
	stderrIn, _ := proc.cmd.StderrPipe()

	log.PushLog("starting %s : %s\n", proc.name, strings.Join(proc.cmd.Args, " "))
	err := proc.cmd.Start()
	if err != nil {
		log.PushLog("error init process %s\n", err.Error())
		return
	}

	go procs.copyAndCapture(proc.name, stdoutIn)
	go procs.copyAndCapture(proc.name, stderrIn)

	go func() {
		proc.cmd.Wait()
		proc.exited = true
	}()
}

func (procs *ChildProcesses) findID(ID ProcessID) bool {
	keys := make([]ProcessID, 0, len(procs.procs))
	for pi := range procs.procs {
		keys = append(keys, pi)
	}

	for _,key := range keys{
		if key == ID {
			return true	
		}
	}
	return false
}


func (procs *ChildProcesses) copyAndCapture(process string, r io.Reader) {
	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf[:])
		if err != nil {
			break
		} else if n ==  0 {
			continue
		}

		lines := strings.Split(string(buf[:n]), "\n")
		for i := 0; i < len(lines); i++ {
			line := lines[i]	
			if len(line) > 0{
				log.PushLog(fmt.Sprintf("Child process (%s): %s",process,line))
			}
		}
	}
}