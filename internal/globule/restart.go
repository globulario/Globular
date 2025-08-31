package globule

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// restart terminates the current Globular process and re-executes it
// with the same binary + args. It preserves environment variables.
// Use only for disruptive config changes (Protocol/Domain).
func (g *Globule) restart() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("restart: resolve executable: %w", err)
	}
	args := os.Args
	env := os.Environ()

	fmt.Println("Restarting Globular...")

	// On Windows, execve semantics differ; just spawn new process and exit.
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

	// On Unix, replace the current process image.
	return syscall.Exec(exe, args, env)
}
