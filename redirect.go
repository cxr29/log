// Copyright (c) 2016 CHEN Xianren. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package log

import (
	"os"
)

func RedirectStdout(name string) error {
	return redirect(name, stdout, os.Stdout)
}

func RedirectStderr(name string) error {
	return redirect(name, stderr, os.Stderr)
}

func redirect(name string, std int, file *os.File) error {
	f, err := openFile(name)
	if err != nil {
		return err
	}

	err = dup2(f, std)
	if err != nil {
		f.Close()
		return err
	}

	*file = *f
	return nil
}
