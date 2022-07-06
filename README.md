# Edit plaintext clipboard as a file

## Why?

As for me, it is very cinvenient way to automate when mixing manual text editing and scripting something.

## Usage

### Launch

    $ go run . -mountpoint some_folder

or, if built

    $ ./clipfs -mountpoint some_folder

and leave it working

* `some_folder` should already exist in Linux and OS X and should not exist in Windows
* launch without parameters to see different options

### Stop

* In Windows, just interrupt the program (Ctrl+C or terminate the process)
* In Unix, just `$ umount some_folder`

### Have fun

It will create virtual FS with files:

* `info` with metainformation about how it works with clipboard
* `clipboard` -- readable (`< clipboard`), writeable (`> clipboard`) and appendable (`>> clipboard`) file, synchronizing with clipboard contents
* `primary` (Linux only) -- readable file wrom which primary (immediate) selection can be read

`clipboard` file is even editable with some editors, which just do write file in-place:  MCEdit, and Vim work correctly.
The most editors fail though: Emacs, Micro, Nano, Far Manager (in Both Windows and Unix) like to create tempoary files, which is not allowed by clipfs.

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



    go run . [-d] -m /tmp/....
