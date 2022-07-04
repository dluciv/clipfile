package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/winfsp/cgofuse/fuse"
)

func main() {
	fDebug := flag.Bool("debug", false, "Debug")
	fCCCP := flag.Bool("cccp", false, "Prefer CCCP backend over everything")
	fWayl := flag.Bool("wayland", false, "Prefer Wayland backend over XClip (when both are available)")
	fMountPoint := flag.String("mountpoint", "", "Mount point")

	flag.Parse()

	if !*fDebug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	if *fMountPoint == "" {
		flag.Usage()
		os.Exit(2)
	}

	if clipboard, err := initClipboards(*fCCCP, *fWayl); err != nil {
		os.Exit(3)
	} else {
		host := fuse.NewFileSystemHost(NewClipFs(clipboard))
		host.Mount(*fMountPoint, []string{})
	}
}
