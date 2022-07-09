# Edit plaintext clipboard as a file

## Why?

IMHO it is very cinvenient way to automate when mixing manual text editing and scripting something.

## Usage

### Launch

    $ go run . -mountpoint some_folder

or, if built

    $ ./clipfs -mountpoint some_folder

and leave it working

* `some_folder` should already exist in Linux and OS X and should not exist in Windows
* launch without parameters to see different options

### Shutdown

* In Windows, just interrupt the program (Ctrl+C or terminate the process)
* In Unix, just `$ umount some_folder`

### Have fun

After launch, it will create virtual FS with files:

* `info` with metainformation about how it works with clipboard
* `clipboard` — readable (`< clipboard`), writeable (`> clipboard`) and appendable (`>> clipboard`) file, synchronizing with clipboard contents
* `primary` (Linux only) — readable (and writable!) file wrom which primary (immediate) selection can be read

`clipboard` and `primary` files is even editable with the most of editors, which just do write file in-place: 
MCEdit, Vim, Emacs, Micro, Nano and Unix-like Far Manager work correctly. Windows File Manager fails right now...

## Build

### Dependencies

* Windows

    > scoop bucket add nonportable
    > scoop install winfsp-np

* Unix-like (Linux, OS X) needs fuse3-compatible library to be installed

### Build

    $ https://github.com/winfsp/cgofuse#how-to-build
    $ go get
    $ go build
