package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Logger struct {
	debugEnable   bool
	app           string
	fileFullName  string
	fileErrorName string
}

var Log Logger

func InitLogger(path, app string, debugEnable bool) {
	fullFileName := fmt.Sprintf("%s/full_%s.log", path, strings.Replace(app, " ", "_", -1))
	errorFileName := fmt.Sprintf("%s/error_%s.log", path, strings.Replace(app, " ", "_", -1))

	Log = Logger{
		debugEnable:   debugEnable,
		app:           app,
		fileFullName:  fullFileName,
		fileErrorName: errorFileName,
	}
}

func writeMsg(paths []string, msg string) {
	var writers []io.Writer

	writers = append(writers, os.Stdout)
	for _, v := range paths {
		f, err := os.OpenFile(v, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		writers = append(writers, f)
	}
	log.SetOutput(io.MultiWriter(writers...))
	log.Print(msg)
}

func formatMsg(lvl, app, msg string) string {
	return fmt.Sprintf("[%s] %s | %s\n", lvl, app, msg)
}

func (l Logger) Fatalf(format string, i ...interface{}) {
	l.Fatal(fmt.Sprintf(format, i...))
}

func (l Logger) Fatal(msg string) {
	fatalMsg := formatMsg("FATAL", l.app, msg)
	writeMsg([]string{l.fileFullName, l.fileErrorName}, fatalMsg)
}

func (l Logger) Warningf(format string, i ...interface{}) {
	l.Warning(fmt.Sprintf(format, i...))
}

func (l Logger) Warning(msg string) {
	warningMsg := formatMsg("WARN", l.app, msg)
	writeMsg([]string{l.fileFullName}, warningMsg)
}

func (l Logger) Infof(format string, i ...interface{}) {
	l.Info(fmt.Sprintf(format, i...))
}

func (l Logger) Info(msg string) {
	infoMsg := formatMsg("INFO", l.app, msg)
	writeMsg([]string{l.fileFullName}, infoMsg)
}

func (l Logger) Debugf(format string, i ...interface{}) {
	l.Debug(fmt.Sprintf(format, i...))
}

func (l Logger) Debug(msg string) {
	if l.debugEnable {
		debugMsg := formatMsg("DEBUG", l.app, msg)
		writeMsg([]string{l.fileFullName}, debugMsg)
	}
}

func (l Logger) Errorf(format string, i ...interface{}) {
	l.Error(fmt.Sprintf(format, i...))
}

func (l Logger) Error(msg string) {
	errorMsg := formatMsg("ERROR", l.app, msg)
	writeMsg([]string{l.fileFullName, l.fileErrorName}, errorMsg)
}
