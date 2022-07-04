package main

import (
	"log"
	"os/exec"

	"github.com/zyedidia/clipper"
)

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

func (c *cccp) WriteAll(reg string, contents []byte) (err error) {
	if reg != "clipboard" {
		return &clipper.ErrInvalidReg{Reg: reg}
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

func initClipboards(preferCCCP, preferWayland bool) (clipper.Clipboard, error) {
	if preferWayland {
		xclipind := -1
		wayldind := -1
		for i, c := range clipper.Clipboards {
			switch c.(type) {
			case *clipper.Xclip:
				xclipind = i
			case *clipper.Wayland:
				wayldind = i
			}
		}
		if wayldind >= xclipind && xclipind >= 0 {
			clipper.Clipboards[xclipind], clipper.Clipboards[wayldind] = clipper.Clipboards[wayldind], clipper.Clipboards[xclipind]
		}
	}
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
