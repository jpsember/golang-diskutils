package snapshot

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/files"
	. "golang-diskutils/gen"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type SnapshotOper struct {
	BaseObject
	config     SnapshotConfig
	outputDir  Path
	imagePaths *Array[Path]
}

func (oper *SnapshotOper) GetHelp(bp *BasePrinter) {
	bp.Pr("take periodic snapshots of screen(s)")
}

func (oper *SnapshotOper) GetArguments() DataClass {
	return DefaultSnapshotConfig
}

func (oper *SnapshotOper) ArgsFileMustExist() bool {
	return false
}

func (oper *SnapshotOper) AcceptArguments(a DataClass) {
	oper.config = a.(SnapshotConfig)
}

func (oper *SnapshotOper) UserCommand() string {
	return "snapshot"
}

func (oper *SnapshotOper) Perform(app *App) {
	oper.SetVerbose(app.Verbose())

	oper.outputDir = NewPathM(oper.config.OutputDir()).GetAbsM()
	if !oper.outputDir.IsDir() {
		BadArg("output directory doesn't exist:", oper.outputDir)
	}

	oper.constructPathBuffer()

	var iter int32

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

func (oper *SnapshotOper) constructPathBuffer() {
	v := NewArray[Path]()
	w := NewDirWalk(oper.outputDir).IncludeExtensions("jpg")
	v.Append(w.FilesRelative()...)
	oper.Log("existing images:", v.Size())
	oper.imagePaths = v
}

func (oper *SnapshotOper) takeSnapshot() {
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
		oper.imagePaths.Add(imagePath)
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

func (oper *SnapshotOper) trimPathBuffer() {
	Todo("have generated scalar int fields default to 'int', not 'int32'")
	for oper.imagePaths.Size() > int(oper.config.MaxImages()) {
		p := oper.outputDir.JoinPathM(oper.imagePaths.First())
		oper.imagePaths.Remove(0, 1)
		oper.Log("Deleting:", p)
		p.DeleteFileM()
	}
}

func (oper *SnapshotOper) getNextImagePath(timestamp int64, deviceNum int) Path {
	return oper.outputDir.JoinM(IntToString(deviceNum) + "_" + strconv.FormatInt(timestamp, 10) + ".jpg")
}
