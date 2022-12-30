package childprocess

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pigeatgarlic/VirtualOS-daemon/log"
)

func (procs *ChildProcesses) handleProcess(id ProcessID) {
	proc := procs.procs[id]

	processname := proc.cmd.Args[0]
	stdoutIn, _ := proc.cmd.StdoutPipe()
	stderrIn, _ := proc.cmd.StderrPipe()

	_log := make([]byte, 0)
	for _, i := range proc.cmd.Args {
		_log = append(_log, append([]byte(i), []byte(" ")...)...)
	}

	log.PushLog("starting %s : %s\n", processname, string(_log))
	err := proc.cmd.Start()
	if err != nil {
		log.PushLog("error init process %s\n", err.Error())
		return
	}

	go procs.copyAndCapture(processname, stdoutIn)
	go procs.copyAndCapture(processname, stderrIn)


	go func() {
		proc.cmd.Wait()
	}()

}

func (procs *ChildProcesses) NewChildProcess(cmd *exec.Cmd) ProcessID {
	if cmd == nil {
		return -1
	}

	procs.mutex.Lock()
	defer func() {
		procs.mutex.Unlock()
		procs.count++
	}()

	id := ProcessID(procs.count)
	procs.procs[id] = &ChildProcess{
		cmd: cmd,
	}

	go func() {
		log.PushLog("process %s, process id %d booting up\n", cmd.Args[0], int(id))
		procs.handleProcess(id)
	}()

	return ProcessID(procs.count)
}

func findChar(dat []byte, find string) (out [][]byte) {
	prev := 0
	for pos, i := range dat {
		if i == []byte(find)[0] {
			out = append(out, dat[prev:pos])
			prev = pos + 1
		}
	}

	out = append(out, dat[prev:])
	return
}

func (procs *ChildProcesses) copyAndCapture(process string, rs ...io.Reader) {
	prefix := []byte(fmt.Sprintf("Child process (%s): ", process))
	after := []byte("")

	buf := make([]byte, 1024)
	var n int
	var err error

	for {
		for _, r := range rs {
			n, err = r.Read(buf[:])
			if err != nil || n != 0 {
				break
			}
		}

		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return
		}

		if n > 0 {
			d := buf[:n]
			lines := findChar(d, "\n")
			for _, line := range lines {
				out := append(prefix, line...)
				out = append(out, after...)

				log.PushLog(string(out))
			}
		}
	}
}