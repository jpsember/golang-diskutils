package internal

import (
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/files"
	"os"
	"time"
)

type ErrLogStruct struct {
	BaseObject
	path        Path
	pathDefined bool
	Errors      int
	Warnings    int
	Clean       bool
	SkipHeader  bool
}

type ErrLog = *ErrLogStruct

func NewErrLog(path Path) ErrLog {
	return &ErrLogStruct{
		path: path,
	}
}

var Warning = Error("Warning")

func (log ErrLog) Add(err error, messages ...any) error {
	errType := "Error"
	if err == Warning {
		errType = "Warning"
		log.Warnings++
	} else {
		log.Errors++
	}
	errMsg := "*** " + errType + ": " + ToString(messages...)
	Pr(errMsg)

	starting := false
	if !log.pathDefined {
		if log.path.NonEmpty() {
			if log.Clean && log.path.Exists() {
				log.path.DeleteFileM()
			}
		}
		log.pathDefined = true
		starting = true
	}
	if log.path.Empty() {
		return nil
	}

	outMsg := errMsg + "\n"

	if starting && !log.SkipHeader {
		outMsg = Dashes + time.Now().Format(time.ANSIC) + "\n" + Dashes + outMsg
	}
	f, err := os.OpenFile(log.path.String(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	CheckOk(err, "Failed opening error file:", log.path)
	defer f.Close()
	_, err = f.WriteString(outMsg)
	CheckOk(err, "Failed appending to error file:", log.path)
	return err
}

func (log ErrLog) PrintSummary() {
	if log.Errors > 0 {
		Pr("*** Total errors:", log.Errors)
	}
	if log.Warnings > 0 {
		Pr("*** Total warnings:", log.Warnings)
	}
	if log.Errors+log.Warnings != 0 {
		Pr("*** See", log.path, "for details.")
	}
}

func (log ErrLog) IssueCount() int {
	return log.Warnings + log.Errors
}
