package main

import (
	"fmt"
	"os"

	"github.com/Divya063/pingApp/cmd"
)

var usage = `
Usage:
    ping [-c count] [-i interval] [-t timeout] [--privileged] host
Examples:
    # ping google continuously
	sudo pingApp ping google.com
    # ping google 5 times
    sudo pingApp ping --count 5 google.com
    # ping google 5 times at 2 seconds intervals
    sudo pingApp ping google.com  --count 3 --interval 2
`

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// p("8.8.8.8")
	// // p("172.27.0.1")
	// p("0:0:0:0:0:ffff:7f00:1")
	// //p("reddit.com")
	// p("2600::")

	//for {
	//    p("google.com")
	//    time.Sleep(1 * time.Second)
	//}
}
