package main

import (
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/json"
	"github.com/jpsember/golang-base/jt"
	"golang-diskutils/gen"
	"testing" // We still need to import the standard testing package
)

var _ = Pr

var tree1 = `
{"a.txt" : "",
 "b.txt" : "",
 "c"     : {"d.txt":"", "e.txt":"", "f" : {"g.txt" : ""}},
}
`

func TestCopyDir(t *testing.T) {
	j := jt.New(t)
	var jsmap = JSMapFromStringM(tree1)

	wd := j.GetTestResultsDir()
	j.GenerateSubdirs(wd.JoinM("source"), jsmap)

	config := gen.NewCopyDirConfig()
	config.SetSource("source")
	config.SetDest("output")
	configPath := wd.JoinM("copydir-args.json")
	configPath.WriteStringM(config.String())

	app := prepareApp()
	oper := &CopyDirOper{}
	oper.ProvideName(oper)
	app.RegisterOper(oper)
	app.SetTestArgs(" -a " + configPath.String() + " --verbose")
	app.Start()

	j.AssertGenerated()
}
