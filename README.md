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
* `clipboard` — readable (`< clipboard`), writeable (`> clipboard`) and even appendable (`>> clipboard`) file, synchronizing with clipboard contents
* `primary` (Linux only) — readable (and writable and appendable!) file to access primary (immediate) selection

`clipboard` and `primary` files is even editable with the most of editors, which just do write file in-place: 
MCEdit, Vim, Emacs, Micro, Nano and Far Manager for Unix work correctly.
Some editors like Far Manager for Windows create temporary files, which clipfs does not allow. The simple workaround for it is
to symlink (for Windows, Far Manager itself can help you =)) clipfile to any regular folder in which it can create temporary files and edit clipboard from there.

## Build

### Dependencies

* Windows

    ```
    > scoop bucket add nonportable
    > scoop install winfsp-np
    ```

* Unix-like (Linux, OS X) needs fuse3-compatible library to be installed

### Build

See [here](https://github.com/winfsp/cgofuse#how-to-build) how to build CGoFuse. Note `CGO_ENABLED` options. Then

    $ go get
    $ go build
