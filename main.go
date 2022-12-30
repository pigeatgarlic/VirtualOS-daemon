package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	childprocess "github.com/pigeatgarlic/VirtualOS-daemon/child-process"
)

var BlacklistedLog = []string{ "adapter.exe" }
var LogDestination = "adapter.exe" 





type Daemon struct {
	childprocess *childprocess.ChildProcesses
	shutdown     chan bool

	domain string
	token  string
}

func TerminateAtTheEnd(daemon *Daemon) {
	chann := make(chan os.Signal, 10)
	signal.Notify(chann, syscall.SIGTERM, os.Interrupt)
	<-chann

	daemon.childprocess.CloseAll()
	time.Sleep(100 * time.Millisecond)
	daemon.shutdown <- true
}

func main() {
	daemon := Daemon{
		shutdown:               make(chan bool),
	}

	args := os.Args[1:]
	for i, arg := range args {
		next := args[i+1]
		switch(arg) {
		case "--domain":
			daemon.domain = next;
		case "--token":
			daemon.token = next;
		default:
		}	
	}



	<-daemon.shutdown
}
