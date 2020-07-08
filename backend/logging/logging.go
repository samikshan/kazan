package logging

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// LevelFlag is the common logging flag to be used to set log level
	LevelFlag = "loglevel"
	// LevelHelp is the help message for LevelFlag
	LevelHelp = "Severity of messages to be logged"

	// YY-MM-DD HH:MM:SS.SSSSSS
	timestampFormat = "2006-01-02 15:04:05.000000"
)

type utcFormatter struct {
	log.Formatter
}

// Format formats the time in UTC timezone
func (u utcFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

func setLogOutput(w io.Writer) {
	log.SetOutput(w)

	// turn off standard library logging
	// see https://github.com/golang/go/issues/19957
	stdlog.SetFlags(0)
	stdlog.SetOutput(ioutil.Discard)
}

type OutputSplitter struct{}

func (splitter *OutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("ERROR")) || bytes.Contains(p, []byte("FATAL")) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

func Init(logLevel string) error {
	l, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		setLogOutput(os.Stderr)
		log.WithError(err).Debug("Failed to parse log level")
		return err
	}
	log.SetLevel(l)
	log.SetOutput(&OutputSplitter{})
	log.SetReportCaller(true)
	log.SetFormatter(utcFormatter{
		&log.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: true,
			FullTimestamp:          true,
			TimestampFormat:        timestampFormat,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				fileName := path.Base(f.File)
				funcName := path.Base(f.Function)

				return fmt.Sprintf("%20s()", funcName), fmt.Sprintf(" %15s:%d", fileName, f.Line)
			},
		},
	})

	return nil
}
