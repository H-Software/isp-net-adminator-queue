package command

import (
	"fmt"

	"github.com/go-cmd/cmd"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
)

type Command struct {
	cmd  string
	args []string
}

var (
	logger  *log.Logger
	command Command
)

const (
	gateway_wifi_fqdn = "10.128.0.2"
	gateway_3_fqdn    = "10.128.0.3"
	ext_scripts_path  = "external_scripts"
)

func init() {
	logger = log.NewLogger(nil)
}

func RunCommand(itemId int) error {

	logger.Infof("running command for ItemId %d", itemId)

	switch itemId {
	case 1:
		// gateway-3 - restriction (net-n/sikana)
		command.cmd = "php"
		command.args = []string{fmt.Sprintf("%v/AdminatorWorkItems/mk_rh_restriction.php", ext_scripts_path), gateway_3_fqdn}

		executeCommand(command.cmd, command.args)

		// $mess_ok = "gateway-3-restriction ok ";
		// $mess_er = "gateway-3-restriction error ";
	case 2, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 33:
		// gateway-wifi (1) - restrictions (net-n/sikana)

		command.cmd = "php"
		command.args = []string{fmt.Sprintf("%v/AdminatorWorkItems/mk_rh_restriction.php", ext_scripts_path), gateway_wifi_fqdn}

		executeCommand(command.cmd, command.args)

		// executeCommand(command.cmd, command.args)

		// $mess_ok = "gateway-wifi-iptables-restart ok ";
		// $mess_er = "gateway-wifi-iptables-restart error ";
	case 3, 4, 5, 6, 7, 8, 9, 10:
		//gateway-fiber (2) - iptables (net-n/sikana)

		// $cmd = "/root/bin/gateway-fiber.remote.exec2.sh \"/etc/init.d/iptables-adminator restart\" ";

		// $mess_ok = "gateway-fiber.iptables ok ";
		// $mess_er = "gateway-fiber.iptables error ";

		command.cmd = "php"
		command.args = []string{"-v"}
		executeCommand(command.cmd, command.args)

	default:
		// unknown itemId
		errM := fmt.Errorf("unsupported ItemId (%d)", itemId)

		logger.Error(errM)
		return fmt.Errorf("%s", errM)
	}

	// TODO: add more commands to case below from original code
	// https://github.com/H-Software/isp-net-adminator/pull/260/files#diff-bf99864cec6493c7e3d8e681dd2fc01c1ffc7480b5728411b20c7af4bbf88b37L874

	return nil
}

func executeCommand(inputCommand string, inputCommandArgs []string) error {

	var errM error

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
				logger.Info(fmt.Printf("COMMAND STDOUT: %v \n",line))
			case line, open := <-cmd.Stderr:
				if !open {
					cmd.Stderr = nil
					continue
				}
				logger.Error(fmt.Printf("COMMAND STDERR: %v \n", line))
			}
		}
	}()

	// Run and wait for Cmd to return
	status := <-cmd.Start()

	// Wait for goroutine to print everything
	<-doneChan

	logger.Debugf("command executed (PID: %v, complete: %v, exit code: %v)",
		status.PID, status.Complete, status.Exit)

	if (status.Exit <= 0 || !status.Complete) {
		errM = fmt.Errorf("command failed! (exitCode: %d, complete: %t)", status.Exit, status.Complete);
	}

	if errM != nil {
		logger.Error(errM)
		return fmt.Errorf("%s", errM)		
	}

	return nil
}
