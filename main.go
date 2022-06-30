package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/winfsp/cgofuse/fuse"
)

func main() {
	fDebug := flag.Bool("d", false, "Debug")
	fMountPoint := flag.String("m", "", "Mount point")

	flag.Parse()

	if *fMountPoint == "" {
		flag.Usage()
		os.Exit(2)
	}

	if !*fDebug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	host := fuse.NewFileSystemHost(&clipFs{})
	host.Mount(*fMountPoint, []string{})
}
