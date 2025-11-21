package cmd

import (
	"bytes"
	"os/exec"
	"time"
)

func ExecWithTimeOut(cmdStr string, timeout time.Duration) (string, error) {
	cmd := exec.Command("sh", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(timeout)
	select {
	case <-after:
		_ = cmd.Process.Kill()
		return "", nil
	case err := <-done:
		if err != nil {
			return "", err
		}
	}
	return stdout.String(), nil
}
