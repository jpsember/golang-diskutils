package names

import (
	. "github.com/jpsember/golang-base/app"
	"github.com/jpsember/golang-base/jt"
	"golang-diskutils/gen"
	. "golang-diskutils/internal"
	"testing"
)

type namesStateStruct struct {
	j      *jt.J
	App    *App
	config gen.NamesConfigBuilder
}

type nstate = *namesStateStruct

func NamesInfo(j *jt.J) nstate {
	info := namesStateStruct{
		j:      j,
		config: gen.NewNamesConfig(),
	}
	return &info
}

func TestNamesWordBackups(t *testing.T) {

	const script = `
{"a.txt" : "",
 "b.txt" : "",
 "c"     : {"d.txt":"", "~$wtf.txt":"", "f" : {"g.txt" : ""}},
}
`
	j := jt.New(t)
	NamesInfo(j).gen(script).execute()
}

func TestNamesStrangeChars(t *testing.T) {

	const script = `
{"a.txt" : "",
 "b.txt" : "",
 "c"     : {"d.txt":"", "à.txt":"", "f" : {"g(par).tàxt" : ""}},
}
`
	j := jt.New(t)
	NamesInfo(j).gen(script).execute()
}

func (t nstate) gen(structure string) nstate {
	var jsmap = JSMapFromStringM(structure)
	t.j.GenerateSubdirs(t.j.GetTestResultsDir().JoinM("source"), jsmap)
	return t
}

func (t nstate) app() *App {
	if t.App == nil {
		t.App = PrepareApp()
		oper := &Struct{}
		oper.ProvideName(oper)
		t.App.RegisterOper(oper)

		config := t.config
		config.SetSource("source")
	}
	return t.App
}

func (t nstate) start() {
	t.app()
	configPath := t.j.GetTestResultsDir().JoinM("names-args.json")
	configPath.WriteStringM(t.config.String())
	t.app().AddTestArgs("-a " + configPath.String())
	if t.j.Verbose() {
		t.app().AddTestArgs("--verbose")
	}
	t.App.Start()
}

func (t nstate) execute() nstate {
	t.start()
	t.j.AssertGenerated()
	return t
}
