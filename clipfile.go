package main

import (
	"log"
	"time"

	"github.com/winfsp/cgofuse/fuse"
	"github.com/zyedidia/clipper"
	"golang.org/x/exp/constraints"
)

// WTF?.. I define it myself?..
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

type pidType = int

type clipFile struct {
	api       clipper.Clipboard
	path      string
	mode      int
	buffer    []byte
	needRead  bool
	needFlush bool
	cTime     fuse.Timespec
	mTime     fuse.Timespec
	aTime     fuse.Timespec
	// TODO:
	// Let them all infere each other for now
	// Later check pid form fuse.Getcontext()
	// rOpenCount uint
	// wOpenCount uint
	// wOpenPid   pidType
	openCount uint
}

func (f *clipFile) size() int {
	if f.openCount == 0 {
		if data, err := f.api.ReadAll(f.path[1:]); err == nil {
			return len(data)
		} else {
			return -fuse.EACCES
		}
	} else if f.mode == fuse.O_WRONLY {
		return 0
	} else {
		return len(f.buffer)
	}
}

func (f *clipFile) open(path string, mode int) int {
	if f.openCount == 0 {
		/*
			fuse.O_RDONLY = 0x0		// <
			fuse.O_WRONLY = 0x1		// >
			fuse.O_RDWR   = 0x2		// >>
			fuse.O_APPEND = 0x400	// >> linux only =(
		*/

		f.needRead = mode != fuse.O_WRONLY
		f.needFlush = mode != fuse.O_RDONLY
		f.mode = mode
		// clipFilesLock.Lock()
		// clipFiles[path] = f
		// clipFilesLock.Unlock()

		// When opening for read, we need to report correct size immediately
		if mode == fuse.O_RDONLY {
			f.read(0)
		}
	} else if f.mode != fuse.O_RDONLY {
		return -fuse.EALREADY
	}
	f.openCount++
	log.Printf(" - - '%s' open count <- %d", f.path, f.openCount)
	return 0
}

func (f *clipFile) read(ofst int64) ([]byte, int) {
	if f.needRead {
		if data, err := f.api.ReadAll(f.path[1:]); err != nil {
			log.Printf(" - - got %d error", err)
			return nil, -1
		} else {
			log.Printf(" - - got '%s' data", data)
			f.buffer = data
			f.needRead = false
		}
	}
	log.Printf(" - - reading clipboard, got '%s'...", string(f.buffer))
	f.aTime = fuse.NewTimespec(time.Now())
	return f.buffer[ofst:], 0
}

func (f *clipFile) write(data []byte, ofst int64) (n int) {
	if f.needRead {
		if _, err := f.read(0); err != 0 {
			return err
		}
	}
	// https://github.com/winfsp/cgofuse/blob/ce7e5a65cac7bacaba0ca95c11610aff8b6e0970/examples/memfs/memfs.go#L301
	endofst := int(ofst) + len(data)
	if endofst > len(f.buffer) {
		f.buffer = append(f.buffer, make([]byte, endofst-len(f.buffer))...)
	}
	n = copy(f.buffer[ofst:endofst], data)
	f.needFlush = true
	/*
		tmsp := fuse.Now()
		node.stat.Ctim = tmsp
		node.stat.Mtim = tmsp
	*/
	return
}

func (f *clipFile) trunc(size int64) int {
	f.needFlush = true
	f.buffer = f.buffer[:min(int(size), len(f.buffer))]
	f.needRead = size != 0
	return 0
}

func (f *clipFile) flush() (err int) {
	if f.needFlush {
		if f.api.WriteAll(f.path[1:], f.buffer) == nil {
			f.needFlush = false
		} else {
			err = -1
		}
	}
	f.mTime = fuse.NewTimespec(time.Now())
	f.aTime = fuse.NewTimespec(time.Now())
	return
}

func (f *clipFile) close() int {
	f.flush()
	// clipFilesLock.Lock()
	// delete(clipFiles, f.path)
	// clipFilesLock.Unlock()
	f.openCount--
	log.Printf(" - - '%s' open count <- %d", f.path, f.openCount)
	if f.openCount == 0 {
		// reset it
		*f = clipFile{
			api:  f.api,
			path: f.path,
		}
	}
	f.aTime = fuse.NewTimespec(time.Now())
	log.Printf(" - - closing '%s'.", f.path)
	return 0
}
