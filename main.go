package main

import (
	"fmt"
	"os"

	"github.com/Divya063/pingApp/cmd"
)

var usage = `
Usage:
sudo pingApp ping host [--count] [--interval] 
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
}
