package main

import (
	"log"
	"os/exec"

	"github.com/zyedidia/clipper"
)

/*
type Clipboard interface {
	// Init initializes the clipboard and returns an error if it is not
	// accessible
	Init() error
	// ReadAll returns the contents of the clipboard register 'reg'
	ReadAll(reg string) ([]byte, error)
	// WriteAll writes 'p' to the clipboard register 'reg'
	WriteAll(reg string, p []byte) error
}
*/

type cccp struct {
	cccpBackend string
}

var _ clipper.Clipboard = (*cccp)(nil)

func (c *cccp) Init() (err error) {
	d, err := exec.Command("cccp", "b").Output()
	c.cccpBackend = string(d)
	return
}

func (c *cccp) ReadAll(reg string) (result []byte, err error) {
	result, err = exec.Command("cccp", "p").Output()
	if err != nil {
		log.Panicf("Failed to paste from clipboard, error: %s", err.Error())
	}
	return
}

type unsupRegister struct {
	register string
}

func (r unsupRegister) Error() string {
	return r.register
}

func (c *cccp) WriteAll(reg string, contents []byte) (err error) {
	if reg != "clipboard" {
		return unsupRegister{register: reg}
	}

	cmd := exec.Command("cccp", "c")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panicf("Failed to copy to clipboard, error: %s", err.Error())
		return
	}

	go func() {
		stdin.Write(contents)
		// io.WriteString(stdin, "\x0d") // ^D
		defer stdin.Close()
	}()

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return
}

var currentClipboard clipper.Clipboard = nil

func initClipboards(preferCCCP bool) (clipper.Clipboard, error) {
	if preferCCCP {
		clipper.Clipboards = append(
			[]clipper.Clipboard{&cccp{}},
			clipper.Clipboards...,
		)
	} else {
		clipper.Clipboards = append(
			clipper.Clipboards,
			[]clipper.Clipboard{&cccp{}}...,
		)
	}
	return clipper.GetClipboard(clipper.Clipboards...)
}
