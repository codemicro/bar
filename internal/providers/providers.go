package providers

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runCommand(program string, args ...string) ([]byte, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.Output()
	if err != nil {
		ne := fmt.Errorf(`failed to execute "%v" (%+v)`, strings.Join(append([]string{program}, args...), " "), err)
		if x, ok := err.(*exec.ExitError); ok {
			return bytes.TrimSpace(x.Stderr), ne
		}
		return nil, ne
	}
	return bytes.TrimSpace(out), err
}
