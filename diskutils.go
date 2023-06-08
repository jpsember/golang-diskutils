package main

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	"golang-diskutils/copydir"
	. "golang-diskutils/internal"
	"golang-diskutils/names"
	"golang-diskutils/snapshot"
)

var _ = Pr

func main() {
	app := PrepareApp()
	addCopyDirOper(app)
	addExamineFilenamesOper(app)
	addSnapshotOper(app)
	app.Start()
}

func addCopyDirOper(app *App) {
	var oper = &copydir.Struct{}
	oper.ProvideName(oper)
	app.RegisterOper(AssertJsonOper(oper))
}

func addExamineFilenamesOper(app *App) {
	var oper = &names.FilenamesOper{}
	oper.ProvideName("names")
	app.RegisterOper(AssertJsonOper(oper))
}

func addSnapshotOper(app *App) {
	Todo("can we eliminate some of this biolerplate?  I.e., have it figure out its own name?")
	var oper = &snapshot.SnapshotOper{}
	oper.ProvideName("snapshot")
	app.RegisterOper(AssertJsonOper(oper))
}
