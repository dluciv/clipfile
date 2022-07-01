package main

import (
	"log"
	"os"

	"github.com/winfsp/cgofuse/fuse"
)

const (
	infoFilename = "info"
	clipFilename = "clip"
)

var infoContents = "?"

type clipFs struct {
	fuse.FileSystemBase
}

func (fs *clipFs) Open(path string, flags int) (errc int, fh uint64) {
	log.Printf("Opening '%s', flags: 0x%x = 0b%b", path, flags, flags)
	switch path {
	case "/" + infoFilename:
		return 0, 0
	case "/" + clipFilename:
		return 0, open(path, flags).fh
	default:
		return -fuse.ENOENT, ^uint64(0)
	}
}

func (fs *clipFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	log.Printf(" - getattr '%s'", path)
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0o555
	case "/" + infoFilename:
		stat.Mode = fuse.S_IFREG | 0o444
		stat.Size = int64(len(infoContents))
	case "/" + clipFilename:
		stat.Mode = fuse.S_IFREG | 0o622
		stat.Size = int64(clipFileSize(path))
		log.Printf(" - - clipfile size: %d", stat.Size)

		// stat.Mode = fuse.S_IFCHR | 0o620 // owner rw, others write
		// stat.Rdev = 0xfefefefefefefefe
	default:
		return -fuse.ENOENT
	}
	stat.Uid = uint32(os.Geteuid())
	stat.Gid = uint32(os.Getegid())
	return 0
}

func (fs *clipFs) Read(path string, buff []byte, ofst int64, fh uint64) int {
	log.Printf(" - read '%s' [%d] @ %d (%d)... ", path, fh, ofst, len(buff))
	switch path {
	case "/" + infoFilename:
		return copy(buff, infoContents[ofst:])
	case "/" + clipFilename:
		dx, _ := getcf(path)
		d, _ := dx.read(ofst)
		log.Printf(" - - read returned '%s'", string(d))
		n := copy(buff, d)
		log.Printf(" - - got: '%s'", string(buff))
		return n
	default:
		return -fuse.ENOENT
	}
}

func (fs *clipFs) Truncate(path string, size int64, fh uint64) int {
	log.Printf(" - truncate '%s' @ %d", path, size)
	f, _ := getcf(path)
	return f.trunc(size)
}

// Write writes data to a file.
// The FileSystemBase implementation returns -ENOSYS.
func (fs *clipFs) Write(path string, buff []byte, ofst int64, fh uint64) int {
	log.Printf(" - write '%s' [%d] @ %d : '%s' ", path, fh, ofst, string(buff))
	f, _ := getcf(path)
	return f.write(buff, ofst)
}

// Flush flushes cached file data.
// The FileSystemBase implementation returns -ENOSYS.
func (fs *clipFs) Flush(path string, fh uint64) int {
	f, _ := getcf(path)

	log.Printf(" - flush '%s'", path)
	return f.flush()
}

// Release closes an open file.
// The FileSystemBase implementation returns -ENOSYS.
func (fs *clipFs) Release(path string, fh uint64) int {
	log.Printf("Close '%s'", path)
	switch path {
	case "/" + infoFilename:
		return 0
	case "/" + clipFilename:
		f, _ := getcf(path)
		f.close()
		return 0
	default:
		return -fuse.ENOENT
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
	return 0
}
