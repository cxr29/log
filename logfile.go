// Copyright (c) 2016 CHEN Xianren. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package log

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Logfile struct {
	mu                               sync.Mutex
	file                             *os.File
	layout, baseName                 string
	size, seconds, maxSize, maxCount int64
}

func Open(baseName string, seconds, maxSize, maxCount int64) (io.WriteCloser, error) {
	return open(".20060102.150405", baseName, seconds, maxSize, maxCount)
}

func OpenDaily(baseName string) (io.WriteCloser, error) {
	return open(".20060102", baseName, 86400, 0, 0)
}

func open(layout, baseName string, seconds, maxSize, maxCount int64) (io.WriteCloser, error) {
	if seconds <= 0 && maxSize <= 0 && maxCount <= 0 {
		return openFile(baseName)
	}
	lf := &Logfile{
		layout:   layout,
		baseName: baseName,
		seconds:  seconds,
		maxSize:  maxSize,
		maxCount: maxCount,
	}
	if err := lf.rotate(); err != nil {
		ErrWarning(lf.Close())
		return nil, err
	}
	go lf.cycle()
	return lf, nil
}

func (lf *Logfile) Close() error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	return lf.replace(nil)
}

func (lf *Logfile) Write(p []byte) (int, error) {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	if lf.file == nil {
		return 0, os.ErrInvalid
	}
	n, err := lf.file.Write(p)
	lf.size += int64(n)
	if err != nil {
		return n, err
	}
	if lf.maxSize > 0 && lf.size > lf.maxSize {
		err = lf.rotate()
	}
	return n, err
}

func (lf *Logfile) replace(newFile *os.File) error {
	if lf.file != nil {
		if err := lf.file.Close(); err != nil {
			return err
		}
	}
	if newFile != nil {
		if stat, err := newFile.Stat(); err != nil {
			return err
		} else {
			lf.size = stat.Size()
		}
	} else {
		lf.size = 0
	}
	lf.file = newFile
	return nil
}

func (lf *Logfile) rotate() error {
	name := lf.baseName + time.Now().Format(lf.layout)
	if lf.file != nil && lf.file.Name() == name {
		return nil
	}
	file, err := openFile(name)
	if err == nil {
		err = lf.replace(file)
	}
	if err != nil {
		return err
	}
	existWarning(os.Remove(lf.baseName))
	_, f := filepath.Split(name)
	ErrWarning(os.Symlink(f, lf.baseName))
	go lf.purge()
	return nil
}

func (lf *Logfile) cycle() {
	if lf.seconds <= 0 {
		return
	}
	ns := lf.seconds * 1e9
	for {
		now := time.Now().UnixNano()
		next := (now/ns)*ns + ns
		<-time.After(time.Duration(next - now))
		lf.mu.Lock()
		if lf.file == nil {
			lf.mu.Unlock()
			return
		} else {
			ErrWarning(lf.rotate())
			lf.mu.Unlock()
		}
	}
}

func (lf *Logfile) purge() {
	if lf.maxCount <= 0 {
		return
	}
	names, _ := filepath.Glob(lf.baseName + "*")
	i, n := 0, len(names)
	for i < n {
		_, err := time.Parse(lf.layout, names[i][len(lf.baseName):])
		if err == nil {
			i++
		} else {
			names = append(names[:i], names[i+1:]...)
			n--
		}
	}
	n -= int(lf.maxCount)
	if n > 0 {
		sort.Strings(names)
		for i = 0; i < n; i++ {
			if lf.file == nil || lf.file.Name() != names[i] {
				existWarning(os.Remove(names[i]))
			}
		}
	}
}

func openFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func existWarning(err error) {
	if err != nil && !os.IsNotExist(err) {
		Warningln(err)
	}
}
