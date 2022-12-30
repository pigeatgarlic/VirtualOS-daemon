package fabric

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pigeatgarlic/VirtualOS-daemon/log"
)

type Command struct {
	Process string `json:"process"`
	Command string `json:"command"`
	Data    string `json:"data"`
}

type Fabric struct {
	Channels map[string]*struct {
		in chan string
		out chan string
	}

	mutex *sync.Mutex
}



func NewFaric() *Fabric{
	ret := Fabric {
		Channels: make(map[string]*struct{in chan string; out chan string}),
		mutex: &sync.Mutex{},
	}


	return &ret;
}



func (f *Fabric)DescribeRoutingRule() {

}
func (f *Fabric)route(source string,cmd Command) (dest string,err error) {
	if source == cmd.Process {
		dest,err = "",fmt.Errorf("loopback command")
	} else {
		dest,err = cmd.Process,nil
	}

	return 
}


func (f *Fabric)HandleNewChannel(process string, in chan string, out chan string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.Channels[process] = &struct{in chan string; out chan string}{
		in: in,
		out: out,
	}	

	go func(_proc string) {
		for {
			f.mutex.Lock()
			chann := f.Channels[_proc];
			f.mutex.Unlock()

			command := Command{}
			err := json.Unmarshal([]byte(<-chann.in),&command)
			if err != nil {
				log.PushLog("error unmarshal message %s",err.Error())
				continue
			}

			dest,err := f.route(_proc,command)
			if err != nil {
				log.PushLog("error routing message %s",err.Error())
				continue
			}


			f.mutex.Lock()
			out_chan := f.Channels[dest].out;
			f.mutex.Unlock()

			dat,_ := json.Marshal(command);
			out_chan <- string(dat)
		}
	}(process)
}