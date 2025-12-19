package globule

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	agentclient "github.com/globulario/Globular/internal/agentclient"
)

// restart terminates the current Globular process and re-executes it
// with the same binary + args. It preserves environment variables.
// Use only for disruptive config changes (Protocol/Domain).
func (g *Globule) restart() error {
	ctx := context.Background()
	if err := agentclient.ApplySingleUnitAction(ctx, NodeAgentAddress(), "globular", "restart"); err == nil {
		return nil
	} else {
		g.log.Warn("node-agent restart failed; falling back to legacy exec", "err", err)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("restart: resolve executable: %w", err)
	}
	args := os.Args
	env := os.Environ()

	fmt.Println("Restarting Globular...")

	if runtime.GOOS == "windows" {
		cmd := exec.Command(exe, args[1:]...)
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("restart (windows): %w", err)
		}
		os.Exit(0)
		return nil
	}

	return syscall.Exec(exe, args, env)
}
