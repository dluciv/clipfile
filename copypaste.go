package main

import (
	"io"
	"log"
	"os/exec"
	"syscall"
)

func clipCopy(contents []byte) syscall.Errno {
	cmd := exec.Command("cccp", "c")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panicf("Failed to copy to clipboard, error code: %d", err)
		return syscall.EIO
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(contents))
		// io.WriteString(stdin, "\x0d") // ^D
	}()

	stdout, err := cmd.CombinedOutput()
	_ = stdout
	if err != nil {
		log.Fatal(err)
		return syscall.EIO
	}

	return syscall.F_OK
}

func clipPaste() ([]byte, syscall.Errno) {
	out, err := exec.Command("cccp", "p").Output()
	if err != nil {
		log.Panicf("Failed to paste from clipboard, error code: %d", err)
		return nil, syscall.EIO
	}
	for _, c := range string(out) {
		log.Printf(" - - - - 0x%x = '%c'", c, c)
	}

	return out, syscall.F_OK
}
