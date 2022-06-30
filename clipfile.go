package main

import (
	"log"
	"sync"

	"github.com/winfsp/cgofuse/fuse"
	"golang.org/x/exp/constraints"
)

// WTF?.. I define it myself?..
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

type clipFile struct {
	fh     uint64 // not recommended to use though
	path   string
	mode   int
	buffer []byte
	//opened    sync.Mutex
	needRead  bool
	needFlush bool
}

var clipFiles = make(map[string]*clipFile)
var clipFilesLock = sync.RWMutex{}

func clipFileSize(path string) int {
	if f, alreadyOpened := getcf(path); !alreadyOpened {
		f = open(path, fuse.O_RDONLY)
		defer f.close()
		return f.size()
	} else {
		return f.size()
	}
}

func (f *clipFile) size() int {
	if f.mode == fuse.O_WRONLY {
		return 0
	} else {
		return len(f.buffer)
	}
}

func getcf(path string) (*clipFile, bool) {
	clipFilesLock.RLock()
	f, ok := clipFiles[path]
	clipFilesLock.RUnlock()
	return f, ok
}

func open(path string, mode int) (f *clipFile) {
	f, exist := getcf(path)
	if !exist {
		/*
			fuse.O_RDONLY = 0x0		// <
			fuse.O_WRONLY = 0x1		// >
			fuse.O_RDWR   = 0x2		// >>
			fuse.O_APPEND = 0x400	// >> linux only =(
		*/
		mode &= fuse.O_RDONLY & fuse.O_WRONLY & fuse.O_RDWR & fuse.O_APPEND
		f = &clipFile{
			fh:        0,
			path:      path,
			mode:      mode,
			needRead:  mode != fuse.O_WRONLY,
			needFlush: mode != fuse.O_RDONLY,
		}
		clipFilesLock.Lock()
		clipFiles[path] = f
		clipFilesLock.Unlock()

		// When opening for read, we need to report correct size immediately
		if mode == fuse.O_RDONLY {
			f.read(0)
		}
	}
	return
}

func (f *clipFile) read(ofst int64) ([]byte, int) {
	if f.needRead {
		if data, err := clipPaste(); err != 0 {
			log.Printf(" - - got %d error", err)
			return nil, int(err)
		} else {
			log.Printf(" - - got '%s' data", data)
			f.buffer = data
			f.needRead = false
		}
	}
	log.Printf(" - - reading clipboard, got '%s'...", string(f.buffer))
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
		err = int(clipCopy(f.buffer))
		if err == 0 {
			f.needFlush = false
		}
	}
	return
}

func (f *clipFile) close() int {
	f.flush()
	clipFilesLock.Lock()
	delete(clipFiles, f.path)
	clipFilesLock.Unlock()
	return 0
}
