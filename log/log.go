package log

import "fmt"

var input  chan string
var output chan string
var configured = false

func init() {
	input = make(chan string, 100)

	go func ()  {
		for {
			out := <-input
			if configured {
				output<-out
			}
		}
	}()
}

func ConfigureDestination(out chan string) {
	output = out
	configured = true
}

func PushLog(format string, a ...any) {
	input<-fmt.Sprintf(format,a...);
}