package command

import (
	"github.com/go-cmd/cmd"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
)

type Command struct {
	cmd string
	args []string
}

var (
	logger *log.Logger
	command Command
)

func init() {
	logger = log.NewLogger(nil)
}

func RunCommand(itemId int) error {
	cmd := prepareCommand(itemId)

	logger.Infof("running command: %s, args: %s", cmd.cmd, cmd.args)

	executeCommand(cmd.cmd, cmd.args)

	return nil
}

func prepareCommand(itemId int) Command {

	command.cmd = "php"
	command.args = []string{"-v"}

	return command
}

func executeCommand(inputCommand string, inputCommandArgs []string) error {

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
				logger.Info(line)
			case line, open := <-cmd.Stderr:
				if !open {
					cmd.Stderr = nil
					continue
				}
				logger.Error(line)
			}
		}
	}()

	// Run and wait for Cmd to return
	status := <-cmd.Start()

	// Wait for goroutine to print everything
	<-doneChan

	logger.Infof("command executed (PID: %v, complete: %v, exit code: %v)",
		status.PID, status.Complete, status.Exit)

	return nil
}
