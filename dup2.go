// Copyright (c) 2016 CHEN Xianren. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build !windows

package log

import (
	"os"
	"syscall"
)

var (
	stdout = syscall.Stdout
	stderr = syscall.Stderr
)

func dup2(f *os.File, fd int) error {
	return syscall.Dup2(int(f.Fd()), fd)
}
