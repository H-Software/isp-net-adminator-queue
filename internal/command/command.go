package command

import (
	"fmt"
	// "regexp"
	"github.com/go-cmd/cmd"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.NewLogger(nil)
}

// func splitBySpaces(str string) []string {
// 	r := regexp.MustCompile("[^\\s]+")
// 	return r.FindAllString(str, -1)
// }

func ExecuteCommand(inputCommand string, inputCommandArgs []string) error {

	logger.Infof("executing command: %s, args: %s", inputCommand, inputCommandArgs)

	// code based on example below
	// https://github.com/go-cmd/cmd/blob/master/examples/blocking-streaming/main.go

	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	// Create Cmd with options
	cmd := cmd.NewCmdOptions(cmdOptions, inputCommand, inputCommandArgs...)

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for cmd.Stdout != nil || cmd.Stderr != nil {
			select {
			case line, open := <-cmd.Stdout:
				if !open {
					cmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case _, open := <-cmd.Stderr:
				if !open {
					cmd.Stderr = nil
					continue
				}
				// fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return
	status := <-cmd.Start()

	// Wait for goroutine to print everything
	<-doneChan

	// logger.Infof("command executed (StdOut: %v, stdErr: %v)", cmd.Stdout, cmd.Stderr)
	logger.Infof("command executed (Status: %v) err: %v stdout: %v, complete: %v, exit code: %v",
		status, cmd.Stderr, cmd.Stdout, status.Complete, status.Exit)

	return nil
}
