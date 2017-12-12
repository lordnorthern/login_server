package server

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"time"
)

// StartCommandConcole will be sitting there and waiting for commands
func StartCommandConcole() {
	SignalChan := make(chan os.Signal)
	signal.Notify(SignalChan, os.Interrupt)

	go func() {
		// Handle the interrupt thing.
		for sig := range SignalChan {
			if sig.String() == "interrupt" {
				processCommands("quit")
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			processCommands(scanner.Text())
		}
	}()
}

func processCommands(str string) {
	switch str {
	case "quit":
		fmt.Println("Terminating services and cleaning up...")
		for _, inst := range Terminatables {
			fmt.Print("Terminating ", inst.GetName(), "...")
			inst.Terminate()
			fmt.Println("Success")
		}
		fmt.Println("Success")
		fmt.Println("Exiting...")
		time.Sleep(time.Millisecond * 500)
		EndServe <- true
	}
}
