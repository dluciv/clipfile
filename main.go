package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/winfsp/cgofuse/fuse"
)

var (
	errLog *log.Logger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	dbgLog *log.Logger = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	fDebug := flag.Bool("debug", false, "Debug")
	fMtime := flag.Bool("mtime", true, "Modify time attributes; switch off if having troubles with your favorite editor")
	fCCCP := flag.Bool("cccp", false, "Prefer CCCP backend over everything")
	fWayl := flag.Bool("wayland", false, "Prefer Wayland backend over XClip (when both are available)")
	fMountPoint := flag.String("mountpoint", "", "Mount point")

	flag.Parse()

	if !*fDebug {
		dbgLog.SetFlags(0)
		dbgLog.SetOutput(ioutil.Discard)
	}

	if *fMountPoint == "" {
		flag.Usage()
		os.Exit(2)
	}

	if clipboard, err := initClipboards(*fCCCP, *fWayl); err != nil {
		os.Exit(3)
	} else {
		host := fuse.NewFileSystemHost(NewClipFs(clipboard, *fMtime))
		host.Mount(*fMountPoint, []string{})
	}
}
