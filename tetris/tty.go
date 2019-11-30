// +build linux darwin
// +build cgo
package main

// #include <termios.h>
// #include <stdio.h>
import "C"
import "fmt"

// TCGetLocalMode get tty local mode setting
func TCGetLocalMode(fd uintptr) (uint64, error) {
	ctermios := C.struct_termios{}
	if e := C.tcgetattr(C.int(fd), &ctermios); e != 0 {
		return 0, fmt.Errorf("%s", C.perror(C.CString("tcgetattr failed")))
	}
	return uint64(ctermios.c_lflag), nil
}

// TCSetLocalMode set tty local mode setting
func TCSetLocalMode(fd uintptr, mode uint64) error {
	ctermios := C.struct_termios{}
	if e := C.tcgetattr(C.int(fd), &ctermios); e != 0 {
		return fmt.Errorf("%s", C.perror(C.CString("tcgetattr failed")))
	}

	ctermios.c_lflag = C.ulong(mode)
	if e := C.tcsetattr(C.int(fd), C.TCSANOW, &ctermios); e != 0 {
		return fmt.Errorf("%s", C.perror(C.CString("tcsetattr failed")))
	}

	return nil
}
