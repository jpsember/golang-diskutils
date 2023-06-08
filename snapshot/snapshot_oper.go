package snapshot

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/files"
	. "golang-diskutils/gen"
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

	oper.outputDir = NewPathM(oper.config.OutputDir())
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
	Todo("takeSnapshot")
}

func (oper *SnapshotOper) trimPathBuffer() {
	Todo("have generated scalar int fields default to 'int', not 'int32'")
	for oper.imagePaths.Size() > int(oper.config.MaxImages()) {
		p := oper.imagePaths.First()
		oper.imagePaths.Remove(0, 1)
		p.DeleteFileM()
	}
}
