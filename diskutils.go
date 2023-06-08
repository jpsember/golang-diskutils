package main

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "golang-diskutils/copydir"
	. "golang-diskutils/internal"
	. "golang-diskutils/names"
)

var _ = Pr

func main() {
	app := PrepareApp()
	addCopyDirOper(app)
	addExamineFilenamesOper(app)
	app.Start()
}

func addCopyDirOper(app *App) {
	var oper = &CopyDirOper{}
	oper.ProvideName(oper)
	app.RegisterOper(AssertJsonOper(oper))
}

func addExamineFilenamesOper(app *App) {
	var oper = &FilenamesOper{}
	oper.ProvideName("names")
	app.RegisterOper(AssertJsonOper(oper))
}
