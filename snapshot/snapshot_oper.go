package snapshot

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "golang-diskutils/gen"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Struct struct {
	BaseObject
	config     SnapshotConfig
	outputDir  Path
	imageNames *Array[string]
}

type Snap = *Struct

func AddOper(app *App) {
	var oper = &Struct{}
	oper.ProvideName("snapshot")
	app.RegisterOper(AssertJsonOper(oper))
}

func (oper Snap) GetHelp(bp *BasePrinter) {
	bp.Pr("take periodic snapshots of screen(s)")
}

func (oper Snap) GetArguments() DataClass {
	return DefaultSnapshotConfig
}

func (oper Snap) ArgsFileMustExist() bool {
	return false
}

func (oper Snap) AcceptArguments(a DataClass) {
	oper.config = a.(SnapshotConfig)
}

func (oper Snap) UserCommand() string {
	return "snapshot"
}

func (oper Snap) Perform(app *App) {
	oper.SetVerbose(app.Verbose())

	oper.outputDir = NewPathM(oper.config.OutputDir()).GetAbsM()
	if !oper.outputDir.IsDir() {
		BadArg("output directory doesn't exist:", oper.outputDir)
	}

	oper.constructPathBuffer()

	var iter int

	for {
		oper.Log("Taking snapshot; iteration", iter)
		oper.takeSnapshot()
		oper.trimPathBuffer()
		time.Sleep(time.Second * time.Duration(oper.config.IntervalSeconds()))
		iter++
		if oper.config.DebugMaxIterations() != 0 {
			if iter >= oper.config.DebugMaxIterations() {
				oper.Log("max iterations reached, quitting")
				break
			}
		}
	}
}

func (oper Snap) constructPathBuffer() {
	v := NewArray[string]()
	w := NewDirWalk(oper.outputDir).IncludeExtensions("jpg")
	for _, x := range w.FilesRelative() {
		v.Add(x.Base())
	}
	oper.Log("existing images:", v.Size())
	oper.imageNames = v
}

func (oper Snap) takeSnapshot() {
	timestamp := time.Now().UnixMilli()
	for devNum := 0; devNum < int(oper.config.NumDevices()); devNum++ {
		imagePath := oper.getNextImagePath(timestamp, devNum)
		s := NewArray[string]()
		dispArg := "-D" + IntToString(1+devNum)
		s.Append(strings.Split("screencapture -S -T 1 -r -tjpg", " ")...)
		s.Add(dispArg)
		s.Add(imagePath.String())
		result, err := makeSysCall(s.Array())
		if strings.Contains(result, "Invalid display specified") {
			Alert("Display "+dispArg, "failed:", INDENT, result)
			continue
		}
		CheckOk(err, "system call failed:", err)
		oper.imageNames.Add(imagePath.Base())
	}
}

func makeSysCall(cmdLineArgs []string) (string, error) {
	var commandName = cmdLineArgs[0]
	var programArguments = cmdLineArgs[1:]
	command := exec.Command(commandName, programArguments...)
	out, err := command.CombinedOutput()
	var strout = string(out)
	return strout, err
}

func (oper Snap) trimPathBuffer() {
	Todo("have generated scalar int fields default to 'int', not 'int32'")
	for oper.imageNames.Size() > int(oper.config.MaxImages()) {
		p := oper.outputDir.JoinM(oper.imageNames.First())
		oper.imageNames.Remove(0, 1)
		oper.Log("Deleting:", p)
		p.DeleteFileM()
	}
}

func (oper Snap) getNextImagePath(timestamp int64, deviceNum int) Path {
	devStr := IntToString(deviceNum)
	timeStr := strconv.FormatInt(timestamp, 10)
	return oper.outputDir.JoinM(timeStr + "_" + devStr + ".jpg")
}
