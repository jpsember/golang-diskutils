package copydir

import (
	. "github.com/jpsember/golang-base/app"
	"github.com/jpsember/golang-base/base"
	"github.com/jpsember/golang-base/jt"
	"golang-diskutils/gen"
	. "golang-diskutils/internal"
	"testing"
)

var tree1 = `
{"a.txt" : "",
 "b.txt" : "",
 "c"     : {"d.txt":"", "e.txt":"", "f" : {"g.txt" : ""}},
}
`

var treeWord = `
{"a.txt" : "",
 "b.txt" : "",
 "c"     : {"d.txt":"", "~$wtf.txt":"", "f" : {"g.txt" : ""}},
}
`

var treeStrange = `
{"a.txt" : "",
 "b.txt" : "",
 "c"     : {"d.txt":"", "à.txt":"", "f" : {"g(par).txt" : ""}},
}
`

type stateStruct struct {
	j      jt.JTest
	App    *App
	config gen.CopyDirConfigBuilder
}

type state = *stateStruct

func NewInfo(j jt.JTest) state {
	info := stateStruct{
		j:      j,
		config: gen.NewCopyDirConfig(),
	}
	return &info
}

func TestCopyDir(t *testing.T) {
	j := jt.New(t)
	NewInfo(j).gen(tree1).execute()
}

func TestWordBackups(t *testing.T) {
	j := jt.New(t)
	NewInfo(j).gen(treeWord).execute()
}

func TestStrangeChars(t *testing.T) {
	j := jt.New(t)
	NewInfo(j).gen(treeStrange).execute()
}

func (t state) gen(structure string) state {
	var jsmap = base.JSMapFromStringM(structure)
	t.j.GenerateSubdirs(t.j.GetTestResultsDir().JoinM("source"), jsmap)
	return t
}

func (t state) app() *App {
	if t.App == nil {
		t.App = PrepareApp()
		oper := &Struct{}
		oper.ProvideName(oper)
		t.App.RegisterOper(oper)

		config := t.config
		config.SetSource("source")
		config.SetDest("output")
	}
	return t.App
}

func (t state) start() {
	t.app()
	configPath := t.j.GetTestResultsDir().JoinM("copydir-args.json")
	configPath.WriteStringM(t.config.String())
	t.app().AddTestArgs("-a " + configPath.String())
	if t.j.Verbose() {
		t.app().AddTestArgs("--verbose")
	}
	t.App.Start()
}

func (t state) execute() state {
	t.start()
	t.j.AssertGenerated()
	return t
}
