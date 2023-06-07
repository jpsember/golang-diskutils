package main

import (
	. "github.com/jpsember/golang-base/app"
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

type UnitTest struct {
	j      *jt.J
	App    *App
	config gen.CopyDirConfigBuilder
}

type UTest = *UnitTest

func NewInfo(t *testing.T) UTest {
	info := UnitTest{
		j:      jt.New(t),
		config: gen.NewCopyDirConfig(),
	}
	return &info
}

func TestCopyDir(t *testing.T) {
	info := NewInfo(t)
	j := info.j

	var jsmap = JSMapFromStringM(tree1)

	j.GenerateSubdirs(j.GetTestResultsDir().JoinM("source"), jsmap)

	info.start()
	j.AssertGenerated()
}

func (t UTest) app() *App {
	if t.App == nil {
		t.App = prepareApp()
		oper := &CopyDirOper{}
		oper.ProvideName(oper)
		t.App.RegisterOper(oper)

		config := t.config
		config.SetSource("source")
		config.SetDest("output")
	}
	return t.App
}

func (t UTest) start() {
	t.app()
	configPath := t.j.GetTestResultsDir().JoinM("copydir-args.json")
	configPath.WriteStringM(t.config.String())
	t.App.AddTestArgs("-a " + configPath.String())
	t.App.AddTestArgs("--verbose")
	t.App.Start()
}
