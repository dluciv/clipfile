package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/winfsp/cgofuse/fuse"
)

func main() {
	fDebug := flag.Bool("debug", false, "Debug")
	fCCCP := flag.Bool("cccp", false, "Prefer CCCP backend")
	fMountPoint := flag.String("mountpoint", "", "Mount point")

	flag.Parse()

	if *fMountPoint == "" {
		flag.Usage()
		os.Exit(2)
	}

	if ccc, err := initClipboards(*fCCCP); err == nil {
		currentClipboard = ccc
	} else {
		os.Exit(3)
	}

	if cc, ok := currentClipboard.(*cccp); ok {
		infoContents = "cccp\n" + cc.cccpBackend + "\n"
	} else {
		infoContents = "clipper\n" + reflect.TypeOf(currentClipboard).Elem().Name() + "\n"
	}

	if !*fDebug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	host := fuse.NewFileSystemHost(&clipFs{})
	host.Mount(*fMountPoint, []string{})
}
