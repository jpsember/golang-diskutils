package main

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	"strings"
)

var _ = Pr

func main() {
	app := prepareApp()
	addCopyDirOper(app)
	addExamineFilenamesOper(app)
	app.Start()
}

func prepareApp() *App {
	var app = NewApp()
	app.SetName("osxutils")
	app.Version = "1.1.0"
	return app
}

const dots = "............................................................................................................................................................................."

func DepthDots(depth int, message ...any) string {
	prefLen := 2 * depth
	CheckState(prefLen < len(dots))
	msg := ToString(message...)
	msg = strings.TrimPrefix(msg, "/")
	return dots[0:prefLen] + " " + msg
}
