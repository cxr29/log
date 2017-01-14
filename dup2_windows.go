// Copyright (c) 2016 CHEN Xianren. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build windows

package log

import (
	"os"
	"syscall"
)

var (
	stdout       = syscall.STD_OUTPUT_HANDLE
	stderr       = syscall.STD_ERROR_HANDLE
	setStdHandle = syscall.NewLazyDLL("kernel32.dll").NewProc("SetStdHandle")
)

func dup2(f *os.File, fd int) error {
	r, _, lastErr := setStdHandle.Call(uintptr(fd), f.Fd())
	if r == 0 {
		return lastErr
	}
	return nil
}
