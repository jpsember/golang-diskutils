package internal

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/files"
)

import "strings"

const dots = "............................................................................................................................................................................."

func DepthDots(depth int, message ...any) string {
	prefLen := 2 * depth
	CheckState(prefLen < len(dots))
	msg := ToString(message...)
	msg = strings.TrimPrefix(msg, "/")
	return dots[0:prefLen] + " " + msg
}

var WindowsTempPattern = Regexp(`^~\$`)

func PrepareApp() *App {
	var app = NewApp()
	app.SetName("osxutils")
	app.Version = "1.1.0"
	return app
}

func RelativePath(path Path, to Path) string {
	pathStr := path.String()
	toStr := to.String()
	i := len(pathStr)
	j := len(toStr)
	if i == 0 || j == 0 || i < j {
		BadArg("can't make:", CR, INDENT, Quoted(pathStr), CR, OUTDENT, "relative to:", CR, INDENT, Quoted(toStr))
	}
	var result string
	if i == j {
		result = ""
	} else {
		result = pathStr[j+1:]
	}
	return result
}

func ProcPath(app *App, desc string, expr string) (Path, string) {
	var err error
	var result Path

	for {
		result, err = NewPath(expr)
		if err != nil {
			break
		}

		result = MakeAbs(result, app.StartDir())
		break
	}
	problem := ""
	if err != nil {
		problem = desc + "; problem: " + err.Error()
	}
	return result, problem
}

func MakeAbs(path Path, absPath Path) Path {
	Todo("put this in the Path package")
	if path.IsAbs() {
		return path
	}
	return absPath.JoinM(path.String())
}
