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
		err = fmt.Errorf(`failed to execute "%v" (%+v)`, strings.Join(append([]string{program}, args...), " "), err)
	}
	return bytes.TrimSpace(out), err
}
