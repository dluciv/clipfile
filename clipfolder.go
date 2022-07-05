package main

import (
	"reflect"
	"time"

	"github.com/winfsp/cgofuse/fuse"
	"github.com/zyedidia/clipper"
)

const (
	infoFilename = "info"
	clipFilename = clipper.RegClipboard
	primFilename = clipper.RegPrimary
)

type clipFs struct {
	fuse.FileSystemBase
	clipboard    *clipFile
	primary      *clipFile
	infoContents string // Small metainfo stored right there
	cTime        fuse.Timespec
}

func NewClipFs(api clipper.Clipboard) (cfs *clipFs) {
	// is there any way to exress this better in Go?..
	hasprimary := false
	if _, ok := api.(*clipper.Xclip); ok {
		hasprimary = true
	} else if _, ok := api.(*clipper.Wayland); ok {
		hasprimary = true
	}

	cfs = &clipFs{}

	if cc, ok := api.(*cccp); ok {
		cfs.infoContents = "API=cccp\nBACKEND=" + cc.cccpBackend + "\n"
	} else {
		cfs.infoContents = "API=clipper\nBACKEND=" + reflect.TypeOf(api).Elem().Name() + "\n"
	}

	cfs.cTime = fuse.NewTimespec(time.Now())

	cfs.clipboard = &clipFile{
		api:   api,
		path:  "/" + clipFilename,
		cTime: cfs.cTime,
		mTime: cfs.cTime,
	}
	if hasprimary {
		cfs.primary = &clipFile{
			api:   api,
			path:  "/" + primFilename,
			cTime: cfs.cTime,
			mTime: cfs.cTime,
		}
	}
	return
}

func (fs *clipFs) getCF(path string) *clipFile {
	switch path {
	case "/" + clipFilename:
		return fs.clipboard
	case "/" + primFilename:
		return fs.primary // might be nil
	default:
		return nil
	}
}

func (fs *clipFs) Open(path string, flags int) (errc int, fh uint64) {
	_, _, pid := fuse.Getcontext() // uid, gid, pid
	dbgLog.Printf("Opening '%s', flags: 0x%x = 0b%b by %d", path, flags, flags, pid)
	switch path {
	case "/" + infoFilename:
		return 0, 0
	default: // clipfiles here
		if cf := fs.getCF(path); cf != nil {
			return cf.open(path, flags), 0
		} else {
			return -fuse.ENOENT, ^uint64(0)
		}
	}
}

func (fs *clipFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	dbgLog.Printf(" - getattr '%s'", path)
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0o555
		stat.Ctim = fs.cTime
		stat.Mtim = fs.cTime
		stat.Atim = fuse.NewTimespec(time.Now())
	case "/" + infoFilename:
		stat.Mode = fuse.S_IFREG | 0o444
		stat.Size = int64(len(fs.infoContents))
		stat.Ctim = fs.cTime
		stat.Mtim = fs.cTime
		stat.Atim = fs.cTime
	default: // clipfiles here
		if cf := fs.getCF(path); cf != nil {
			stat.Mode = fuse.S_IFREG | 0o622
			stat.Size = int64(cf.size())

			stat.Ctim = cf.cTime
			stat.Mtim = cf.mTime
		} else {
			return -fuse.ENOENT
		}
		dbgLog.Printf(" - - clipfile size: %d", stat.Size)
	}
	uid, gid, _ := fuse.Getcontext() // uid, gid, pid
	stat.Uid = uid
	stat.Gid = gid
	return 0
}

func (fs *clipFs) Read(path string, buff []byte, ofst int64, fh uint64) int {
	dbgLog.Printf(" - read '%s' [%d] @ %d (%d)... ", path, fh, ofst, len(buff))
	switch path {
	case "/" + infoFilename:
		return copy(buff, fs.infoContents[ofst:])
	default: // clipfiles here
		if cf := fs.getCF(path); cf != nil {
			d, _ := cf.read(ofst)
			dbgLog.Printf(" - - read returned '%s'", string(d))
			n := copy(buff, d)
			dbgLog.Printf(" - - got: '%s'", string(buff))
			return n
		} else {
			return -fuse.ENOENT
		}
	}
}

func (fs *clipFs) Truncate(path string, size int64, fh uint64) int {
	dbgLog.Printf(" - truncate '%s' @ %d", path, size)

	if cf := fs.getCF(path); cf != nil {
		return cf.trunc(size)
	} else {
		return -fuse.ENOENT
	}
}

// Write writes data to a file.
// The FileSystemBase implementation returns -ENOSYS.
func (fs *clipFs) Write(path string, buff []byte, ofst int64, fh uint64) int {
	dbgLog.Printf(" - write '%s' [%d] @ %d : '%s' ", path, fh, ofst, string(buff))

	if cf := fs.getCF(path); cf != nil {
		return cf.write(buff, ofst)
	} else {
		return -fuse.ENOENT
	}
}

// Flush flushes cached file data.
// The FileSystemBase implementation returns -ENOSYS.
func (fs *clipFs) Flush(path string, fh uint64) int {
	dbgLog.Printf(" - flush '%s'", path)

	if cf := fs.getCF(path); cf != nil {
		return cf.flush()
	} else {
		return -fuse.ENOENT
	}
}

// Release closes an open file.
// The FileSystemBase implementation returns -ENOSYS.
func (fs *clipFs) Release(path string, fh uint64) int {
	_, _, pid := fuse.Getcontext() // uid, gid, pid
	dbgLog.Printf("Close '%s' by %d", path, pid)
	switch path {
	case "/" + infoFilename:
		return 0
	default: // clipfiles here
		if cf := fs.getCF(path); cf != nil {
			return cf.close()
		} else {
			return -fuse.ENOENT
		}

	}
}

func (fs *clipFs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	fill(infoFilename, nil, 0)
	fill(clipFilename, nil, 0)
	if fs.primary != nil {
		fill(primFilename, nil, 0)
	}
	return 0
}
