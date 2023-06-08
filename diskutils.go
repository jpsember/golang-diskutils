package main

import (
	. "github.com/jpsember/golang-base/base"
	"golang-diskutils/copydir"
	. "golang-diskutils/internal"
	"golang-diskutils/names"
	"golang-diskutils/snapshot"
)

var _ = Pr

func main() {
	app := PrepareApp()
	copydir.AddOper(app)
	names.AddOper(app)
	snapshot.AddOper(app)
	app.Start()
}
