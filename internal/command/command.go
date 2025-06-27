package command

import (
	"fmt"
	"net"
	"github.com/go-cmd/cmd"
	"github.com/h-software/isp-net-adminator-queue/internal/log"
	"github.com/h-software/isp-net-adminator-queue/internal/flag"
)

type Command struct {
	cmd  string
	args []string
}

const (
	ext_scripts_path  = "external_scripts"
)

var (
	logger          *log.Logger
	command         Command
	err             error
	work_items_path = fmt.Sprintf("%s/AdminatorWorkItems/src", ext_scripts_path)
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
		command.args = []string{fmt.Sprintf("%v/mk_rh_restriction.php", work_items_path), getIpAddress(flag.ConfigGateway3Fqdn)}

		err = ExecuteCommand(command.cmd, command.args)

		// $mess_ok = "gateway-3-restriction ok ";
		// $mess_er = "gateway-3-restriction error ";
	case 2:
		// gateway-wifi (1) - restrictions (net-n/sikana)
		command.cmd = "php"
		command.args = []string{fmt.Sprintf("%v/mk_rh_restriction.php", work_items_path), getIpAddress(flag.ConfigGatewayWifiFqdn)}
		
		err = ExecuteCommand(command.cmd, command.args)

		// $mess_ok = "gateway-wifi-iptables-restart ok ";
		// $mess_er = "gateway-wifi-iptables-restart error ";
	case 3:
		//gateway-fiber (2) - iptables (net-n/sikana)

		// $cmd = "/root/bin/gateway-fiber.remote.exec2.sh \"/etc/init.d/iptables-adminator restart\" ";

		// $mess_ok = "gateway-fiber.iptables ok ";
		// $mess_er = "gateway-fiber.iptables error ";
	case 4:
		//gateway-fiber - radius
        // $cmd = "/root/bin/reinhard-fiber.remote.exec2.sh \"/root/bin/radius.restart.sh\"";

        // $mess_ok = "reinhard-fiber.radius ok ";
        // $mess_er = "reinhard-fiber.radius error ";
	case 5:
		// $cmd = "/root/bin/gateway-fiber.remote.exec2.sh \"/etc/init.d/shaper restart\" ";

        // $mess_ok = "reinhard-fiber.shaper ok ";
        // $mess_er = "reinhard-fiber.shaper error ";
	case 6:
        // $cmd = "/root/bin/trinity.local.exec2.sh \"/root/bin/mikrotik.dhcp.leases.erase.sh\" ";

        // $mess_ok = "(trinity) mikrotik.dhcp.leases.erase ok ";
        // $mess_er = "(trinity) mikrotik.dhcp.leases.erase error ";

	case 7:
		// $cmd = "/root/bin/trinity.local.exec2.sh \"/root/bin/scripts_fiber/sw.h3c.vlan.set.pl update\" ";

        // $mess_ok = "trinity.sw.h3c.vlan.set ok ";
        // $mess_er = "trinity.sw.h3c.vlan.set error ";

	case 8:
		// nothing
	
	case 9:
        // $cmd = "/root/bin/erik.remote.exec.sh \"/root/bin/dns.restart.sh\" ";

        // $mess_ok = "erik-dns.restart ok ";
        // $mess_er = "erik-dns.restart-restart error ";
	case 10:
		// $cmd = "/root/bin/trinity.local.exec2.sh \"/root/bin/dns.restart.sh\" ";

        // $mess_ok = "trinity-dns-restart ok ";
        // $mess_er = "trinity-dns-restart error ";
	case 11:
		// $cmd = "/root/bin/artemis.remote.exec2.sh \"/root/bin/dns.restart.sh\" ";

        // $mess_ok = "artemis-dns-server-restart ok ";
        // $mess_er = "artemis-dns-server-restart error ";
	case 12:
	    // $cmd = "/root/bin/c.ns.remote.exec2.sh \"/root/bin/dns.restart.sh\" ";

        // $mess_ok = "c.ns.simelon.net-dns-server-restart ok ";
        // $mess_er = "c.ns.simelon.net-dns-server-restart error ";
	case 13:
		// gateway-wifi (ros) - shaper (client's tariffs)

		// "/root/bin/trinity.local.exec2.sh \"php /var/www/html/htdocs.ssl/adminator2/mk_control/mk_qos_handler.php 10.128.0.2\" ";

		// $mess_ok = "gateway-wifi-shaper-restart ok ";
		// $mess_er = "gateway-wifi-shaper-restart error ";

		command.cmd = "php"
		command.args = []string{fmt.Sprintf("%v/mk_qos_handler.php", work_items_path), getIpAddress(flag.ConfigGatewayWifiFqdn)}
		err = ExecuteCommand(command.cmd, command.args)
	case 14:
		// $cmd = "/root/bin/trinity.local.exec2.sh \"/root/bin/scripts_wifi_network/rb.filter_v2.pl\" ";

		// $mess_ok = "trinity-filtrace-IP-on-Mtik's-restart ok ";
        // $mess_er = "trinity-filtrace-IP-on-Mtik's-restart error ";

	case 15:
		//trinity - Monitoring I - Footer-restart (alarms)

		// $cmd = "/root/bin/monitoring.remote.exec2.sh \"/var/www/cgi-bin/mon1-footer.pl\" ";

        // $mess_ok = "monitoring-I-Footer-restart ok ";
        // $mess_er = "monitoring-I-Footer-restart error ";

	case 16:
		// $cmd = "/root/bin/trinity.local.exec2.sh \"/var/www/cgi-bin/cgi-mon/footer_php.pl\" ";

        // $mess_ok = "trinity-monitoring-I-Footer-PHP-restart ok ";
        // $mess_er = "trinity-monitoring-I-Footer-PHP-restart error ";

	case 17:
        // $cmd = "/root/bin/trinity.local.exec2.sh \"/var/www/cgi-bin/cgi-mon/footer_cat.pl\" ";

        // $mess_ok = "trinity-monitoring-I-Footer-cat-restart ok ";
        // $mess_er = "trinity-monitoring-I-Footer-cat-restart error ";
	case 18:
		// $cmd = "/root/bin/monitoring.remote.exec2.sh \"/var/www/cgi-bin/mon2-feeder.pl\" ";

        // $mess_ok = "monitoring - Monitoring II - Feeder-restart ok ";
        // $mess_er = "monitoring - Monitoring II - Feeder-restart error ";
	case 19:
		// nothing
	case 20:
        // $cmd = "/root/bin/trinity.local.exec2.sh \"php /var/www/html/htdocs.ssl/adminator2/mk_control/mk_qos_handler.php 10.128.0.3\" ";

        // $mess_ok = "gateway-3 (ros) - shaper (client's tariffs) - restart ok ";
        // $mess_er = "gateway-3 (ros) - shaper (client's tariffs) - restart error ";
		
		command.cmd = "php"
		command.args = []string{fmt.Sprintf("%v/mk_qos_handler.php", work_items_path), getIpAddress(flag.ConfigGateway3Fqdn)}
		err = ExecuteCommand(command.cmd, command.args)

	case 21:
		// $cmd = "/root/bin/artemis.remote.exec2.sh \"/root/bin/radius.restart.sh\" ";

        // $mess_ok = "artemis-radius-restart ok ";
        // $mess_er = "artemis-radius-restart error ";
	case 22:
		// $cmd = "/root/bin/monitoring.remote.exec2.sh \"/var/www/cgi-bin/mon2-checker.pl\" ";

        // $mess_ok = "monitoring - Monitoring II - Feeder-restart ok ";
        // $mess_er = "monitoring - Monitoring II - Feeder-restart error ";
	case 23:
		// $cmd = "/root/bin/trinity.local.exec2.sh \"php /var/www/html/htdocs.ssl/adminator2/mk_control/mk_qos_handler.php 10.128.0.15\" ";

        // $mess_ok = "gateway-5-shaper-restart ok ";
        // $mess_er = "gateway-5-shaper-restart error ";
	case 24:
		// $cmd = "/root/bin/trinity.local.exec2.sh \"php /var/www/html/htdocs.ssl/adminator2/mk_control/mk_rh_restriction.php 10.128.0.15\" ";

        // $mess_ok = "gateway-5-iptables-restart ok ";
        // $mess_er = "gateway-5-iptables-restart error ";
	default:
		// unknown itemId
		errM := fmt.Errorf("unsupported ItemId (%d)", itemId)

		logger.Error(errM)
		return fmt.Errorf("%w", errM)
	}

	if err != nil {
		errM := fmt.Sprintf("runCommand failed: %v, (itemId: %d)", err, itemId)
		logger.Error(errM)
		return err
	}

	return nil
}

func ExecuteCommand(inputCommand string, inputCommandArgs []string) error {

	var errM error

	logger.Infof("executing command: %s, args: %s", inputCommand, inputCommandArgs)

	// code based on example below
	// https://github.com/go-cmd/cmd/blob/master/examples/blocking-streaming/main.go

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
				logger.Infof("COMMAND STDOUT: %s", line)
			case line, open := <-cmd.Stderr:
				if !open {
					cmd.Stderr = nil
					continue
				}
				logger.Errorf("COMMAND STDERR: %s", line)
			}
		}
	}()

	// Run and wait for Cmd to return
	status := <-cmd.Start()

	// Wait for goroutine to print everything
	<-doneChan

	logger.Debugf("command executed (PID: %v, complete: %v, exit code: %v)",
		status.PID, status.Complete, status.Exit)

	if status.Exit != 0 || !status.Complete {
		errM = fmt.Errorf("command failed! (exitCode: %d, complete: %t)", status.Exit, status.Complete)
	}

	if errM != nil {
		logger.Error(errM)
		return fmt.Errorf("%w", errM)
	}

	return nil
}

func getIpAddress (fqdn string) string {
	addr, err := net.ResolveIPAddr("ip", fqdn)
	if err != nil {
		panic(err)
	}
	return addr.String()
}
