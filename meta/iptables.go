package meta

import (
	"context"
	"net"
	"os/exec"

	"github.com/misha-ridge/x/tcontext"
	"github.com/misha-ridge/x/tlog"
	"go.uber.org/zap"
)

// ConfigureIPTables configures iptables
func ConfigureIPTables(ctx context.Context, allowedNetwork net.IPNet, addr string) error {
	params := []string{
		"-w", "-t", "nat", "-p", "tcp", "-d", IP, "--dport", "80", "-s", allowedNetwork.String(),
		"-j", "DNAT", "--to-destination", addr,
	}
	// redirect traffic to meta server
	cmd := exec.CommandContext(ctx, "iptables", append([]string{"-I", "PREROUTING", "1"}, params...)...)
	if err := cmd.Run(); err != nil {
		return err
	}
	defer func() {
		// remove iptables rule
		err := exec.CommandContext(tcontext.Reopen(ctx), "iptables", append([]string{"-D", "PREROUTING"}, params...)...).Run()
		if err != nil {
			tlog.Get(ctx).Error("Failed to remove iptables rule", zap.Error(err))
		}
	}()
	<-ctx.Done()
	return ctx.Err()
}
