// Copyright (c) 2016 CHEN Xianren. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = log.LstdFlags
)

const (
	LevelDebug = (iota + 1) * 100
	LevelInfo
	LevelNotice
	LevelWarning
	LevelError
	LevelCritical
	LevelPanic
	LevelFatal
)

var (
	levels = map[int]string{
		LevelDebug:    "DEBUG",
		LevelInfo:     "INFO",
		LevelNotice:   "NOTICE",
		LevelWarning:  "WARNING",
		LevelError:    "ERROR",
		LevelCritical: "CRITICAL",
		LevelPanic:    "PANIC",
		LevelFatal:    "FATAL",
	}
)

const (
	namePrefix = "LEVEL"
	levelDepth = 4
)

func AddBracket() {
	for k, v := range levels {
		levels[k] = "[" + v + "]"
	}
}

func AddColon() {
	for k, v := range levels {
		levels[k] = v + ":"
	}
}

func SetLevelName(level int, name string) {
	levels[level] = name
}

func LevelName(level int) string {
	if name, ok := levels[level]; ok {
		return name
	}
	return namePrefix + strconv.Itoa(level)
}

func NameLevel(name string) int {
	for k, v := range levels {
		if v == name {
			return k
		}
	}
	var level int
	if strings.HasPrefix(name, namePrefix) {
		level, _ = strconv.Atoi(name[len(namePrefix):])
	}
	return level
}

type Logger struct {
	level  int
	logger *log.Logger
}

func New(out io.Writer, prefix string, flag, level int) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(out, prefix, flag),
	}
}

func (l *Logger) Flags() int {
	return l.logger.Flags()
}

func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) Prefix() string {
	return l.logger.Prefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Logger) Level() int {
	return l.level
}

// SetLevel is not locked.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) output(level, calldepth int, s string) error {
	if l == std {
		calldepth++
	}
	return l.logger.Output(calldepth, LevelName(level)+" "+s)
}

func (l *Logger) Err(level, calldepth int, err error) error {
	if err != nil && level >= l.level {
		return l.output(level, calldepth, err.Error())
	}
	return nil
}

func (l *Logger) Output(level, calldepth int, a ...interface{}) error {
	if level >= l.level {
		return l.output(level, calldepth, fmt.Sprint(a...))
	}
	return nil
}

func (l *Logger) Outputf(level, calldepth int, format string, a ...interface{}) error {
	if level >= l.level {
		return l.output(level, calldepth, fmt.Sprintf(format, a...))
	}
	return nil
}

func (l *Logger) Outputln(level, calldepth int, a ...interface{}) error {
	if level >= l.level {
		return l.output(level, calldepth, fmt.Sprintln(a...))
	}
	return nil
}

func (l *Logger) ErrDebug(err error) {
	l.Err(LevelDebug, levelDepth, err)
}

func (l *Logger) ErrNotice(err error) {
	l.Err(LevelNotice, levelDepth, err)
}

func (l *Logger) ErrInfo(err error) {
	l.Err(LevelInfo, levelDepth, err)
}

func (l *Logger) ErrWarning(err error) {
	l.Err(LevelWarning, levelDepth, err)
}

func (l *Logger) ErrError(err error) {
	l.Err(LevelError, levelDepth, err)
}

func (l *Logger) ErrCritical(err error) {
	l.Err(LevelCritical, levelDepth, err)
}

func (l *Logger) ErrPanic(err error) {
	if err != nil {
		l.Err(LevelPanic, levelDepth, err)
		panic(err)
	}
}

func (l *Logger) ErrFatal(err error) {
	if err != nil {
		l.Err(LevelFatal, levelDepth, err)
		os.Exit(1)
	}
}

func (l *Logger) Debug(a ...interface{}) {
	l.Output(LevelDebug, levelDepth, a...)
}

func (l *Logger) Notice(a ...interface{}) {
	l.Output(LevelNotice, levelDepth, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.Output(LevelInfo, levelDepth, a...)
}

func (l *Logger) Warning(a ...interface{}) {
	l.Output(LevelWarning, levelDepth, a...)
}

func (l *Logger) Error(a ...interface{}) {
	l.Output(LevelError, levelDepth, a...)
}

func (l *Logger) Critical(a ...interface{}) {
	l.Output(LevelCritical, levelDepth, a...)
}

func (l *Logger) Panic(a ...interface{}) {
	s := fmt.Sprint(a...)
	if LevelPanic >= l.level {
		l.output(LevelPanic, levelDepth-1, s)
	}
	panic(s)
}

func (l *Logger) Fatal(a ...interface{}) {
	l.Output(LevelFatal, levelDepth, a...)
	os.Exit(1)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Outputf(LevelDebug, levelDepth, format, a...)
}

func (l *Logger) Noticef(format string, a ...interface{}) {
	l.Outputf(LevelNotice, levelDepth, format, a...)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.Outputf(LevelInfo, levelDepth, format, a...)
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Outputf(LevelWarning, levelDepth, format, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Outputf(LevelError, levelDepth, format, a...)
}

func (l *Logger) Criticalf(format string, a ...interface{}) {
	l.Outputf(LevelCritical, levelDepth, format, a...)
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	if LevelPanic >= l.level {
		l.output(LevelPanic, levelDepth-1, s)
	}
	panic(s)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Outputf(LevelFatal, levelDepth, format, a...)
	os.Exit(1)
}

func (l *Logger) Debugln(a ...interface{}) {
	l.Outputln(LevelDebug, levelDepth, a...)
}

func (l *Logger) Infoln(a ...interface{}) {
	l.Outputln(LevelInfo, levelDepth, a...)
}

func (l *Logger) Noticeln(a ...interface{}) {
	l.Outputln(LevelNotice, levelDepth, a...)
}

func (l *Logger) Warningln(a ...interface{}) {
	l.Outputln(LevelWarning, levelDepth, a...)
}

func (l *Logger) Errorln(a ...interface{}) {
	l.Outputln(LevelError, levelDepth, a...)
}

func (l *Logger) Criticalln(a ...interface{}) {
	l.Outputln(LevelCritical, levelDepth, a...)
}

func (l *Logger) Panicln(a ...interface{}) {
	s := fmt.Sprintln(a...)
	if LevelPanic >= l.level {
		l.output(LevelPanic, levelDepth-1, s)
	}
	panic(s)
}

func (l *Logger) Fatalln(a ...interface{}) {
	l.Outputln(LevelFatal, levelDepth, a...)
	os.Exit(1)
}

var std = New(os.Stderr, "", LstdFlags, LevelInfo)

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func Flags() int {
	return std.Flags()
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func Prefix() string {
	return std.Prefix()
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func Level() int {
	return std.Level()
}

// SetLevel is not locked.
func SetLevel(level int) {
	std.SetLevel(level)
}

func Err(level, calldepth int, err error) error {
	return std.Err(level, calldepth, err)
}

func Output(level, calldepth int, a ...interface{}) error {
	return std.Output(level, calldepth, a...)
}

func Outputf(level, calldepth int, format string, a ...interface{}) error {
	return std.Outputf(level, calldepth, format, a...)
}

func Outputln(level, calldepth int, a ...interface{}) error {
	return std.Outputln(level, calldepth, a...)
}

func ErrDebug(err error) {
	std.ErrDebug(err)
}

func ErrInfo(err error) {
	std.ErrInfo(err)
}

func ErrNotice(err error) {
	std.ErrNotice(err)
}

func ErrWarning(err error) {
	std.ErrWarning(err)
}

func ErrError(err error) {
	std.ErrError(err)
}

func ErrCritical(err error) {
	std.ErrCritical(err)
}

func ErrPanic(err error) {
	std.ErrPanic(err)
}

func ErrFatal(err error) {
	std.ErrFatal(err)
}

func Debug(a ...interface{}) {
	std.Debug(a...)
}

func Info(a ...interface{}) {
	std.Info(a...)
}

func Notice(a ...interface{}) {
	std.Notice(a...)
}

func Warning(a ...interface{}) {
	std.Warning(a...)
}

func Error(a ...interface{}) {
	std.Error(a...)
}

func Critical(a ...interface{}) {
	std.Critical(a...)
}

func Panic(a ...interface{}) {
	std.Panic(a...)
}

func Fatal(a ...interface{}) {
	std.Fatal(a...)
}

func Debugf(format string, a ...interface{}) {
	std.Debugf(format, a...)
}

func Infof(format string, a ...interface{}) {
	std.Infof(format, a...)
}

func Noticef(format string, a ...interface{}) {
	std.Noticef(format, a...)
}

func Warningf(format string, a ...interface{}) {
	std.Warningf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	std.Errorf(format, a...)
}

func Criticalf(format string, a ...interface{}) {
	std.Criticalf(format, a...)
}

func Panicf(format string, a ...interface{}) {
	std.Panicf(format, a...)
}

func Fatalf(format string, a ...interface{}) {
	std.Fatalf(format, a...)
}

func Debugln(a ...interface{}) {
	std.Debugln(a...)
}

func Infoln(a ...interface{}) {
	std.Infoln(a...)
}

func Noticeln(a ...interface{}) {
	std.Noticeln(a...)
}

func Warningln(a ...interface{}) {
	std.Warningln(a...)
}

func Errorln(a ...interface{}) {
	std.Errorln(a...)
}

func Criticalln(a ...interface{}) {
	std.Criticalln(a...)
}

func Panicln(a ...interface{}) {
	std.Panicln(a...)
}

func Fatalln(a ...interface{}) {
	std.Fatalln(a...)
}
