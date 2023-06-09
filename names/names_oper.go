package names

import (
	"fmt"
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/files"
	. "golang-diskutils/gen"
	. "golang-diskutils/internal"
	"os"
	"regexp"
	"strings"
)

type Struct struct {
	BaseObject
	errLog          ErrLog
	errPath         Path
	sourcePath      Path
	config          NamesConfig
	pattern         *regexp.Regexp
	deleteFlag      bool
	sourcePrefixLen int
	quitting        bool
}

type Names = *Struct

func AddOper(app *App) {
	var oper = &Struct{}
	oper.ProvideName("names")
	app.RegisterOper(AssertJsonOper(oper))
}

type DirInfoStruct struct {
	Path      Path
	Depth     int
	DiskUsage int64
	Children  *Array[DirInfo]
}

type DirInfo = *DirInfoStruct

func NewDirInfo(path Path, parent *DirInfoStruct) DirInfo {
	s := DirInfoStruct{
		Path:     path,
		Depth:    0,
		Children: NewArray[DirInfo](),
	}
	if parent != nil {
		s.Depth = parent.Depth + 1
	}
	return &s
}

func (oper Names) GetArguments() DataClass {
	return DefaultNamesConfig
}

func (oper Names) ArgsFileMustExist() bool {
	return false
}

func (oper Names) AcceptArguments(a DataClass) {
	oper.config = a.(NamesConfig)
}

func (oper Names) UserCommand() string {
	return "names"
}

func (oper Names) relToSource(path Path) string {
	return RelativePath(path, oper.sourcePath)
}

func (oper Names) Perform(app *App) {
	Todo("Option to rename files, e.g. trimming whitespace, changing dashes")
	oper.SetVerbose(app.Verbose())
	oper.pattern = Regexp(oper.config.Pattern())
	{
		var operSourceDir Path
		problem := ""
		for {
			operSourceDir, problem = ProcPath(app, "Source directory", oper.config.Source())
			if problem != "" {
				break
			}
			if !operSourceDir.IsDir() {
				problem = "source is not a directory: " + operSourceDir.String()
				break
			}
			break
		}
		if problem != "" {
			Pr("Problem:", problem)
			os.Exit(1)
		}
		oper.sourcePath = operSourceDir
	}

	oper.sourcePrefixLen = len(oper.sourcePath.Parent().String()) + 1 // add 1 for separator

	logPath := NewPathOrEmptyM(oper.config.Log())
	if logPath.NonEmpty() {
		logPath = app.StartDir().JoinM(logPath.String())
		Todo("joined logPath to startDir:" + app.StartDir().String() + " result: " + logPath.String())
	}
	Todo("attempting to open logPath: " + logPath.String())
	oper.errLog = NewErrLog(logPath)
	Todo("it was opened")
	oper.errLog.SkipHeader = app.HasTestArgs()
	oper.errLog.Clean = oper.config.CleanLog()

	rootInfo := NewDirInfo(oper.sourcePath, nil)

	oper.processDir(rootInfo)
	Pr(Dashes)
	oper.printDiskUsage(rootInfo)
	Pr(Dashes)
	oper.errLog.PrintSummary()

	Todo("and we are done performing")
}

func (oper Names) processDir(dirInfo DirInfo) {
	if oper.quitting {
		return
	}

	maxIssues := oper.config.MaxProblems()
	if maxIssues > 0 && oper.errLog.IssueCount() >= int(maxIssues) {
		oper.quitting = true
		oper.errLog.Add(Warning, "Stopping since max issue count has been reached")
		return
	}

	dir := dirInfo.Path

	oper.examineFilename(dir)
	if oper.processDeleteFlag(dir) {
		return
	}
	dirEntries, err := os.ReadDir(dir.String())
	if err != nil {
		oper.errLog.Add(err, "unable to ReadDir", oper.relToSource(dir))
		return
	}

	for _, dirEntry := range dirEntries {
		nm := dirEntry.Name()

		sourceFile := dir.JoinM(nm)
		oper.examineFilename(sourceFile)
		if oper.processDeleteFlag(sourceFile) {
			continue
		}

		// Check if source is a symlink.  If so, skip it.
		srcFileInfo, err := os.Lstat(sourceFile.String())
		if err != nil {
			oper.errLog.Add(err, "unable to Lstat", sourceFile)
			continue
		}
		if srcFileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			oper.errLog.Add(Warning, "Found symlink:", sourceFile)
			continue
		}

		if sourceFile.IsDir() {
			child := NewDirInfo(sourceFile, dirInfo)
			dirInfo.Children.Add(child)
			oper.processDir(child)
			continue
		}

		oper.Log(DepthDots(dirInfo.Depth, RelativePath(sourceFile, oper.sourcePath)))

		stat, err := os.Stat(sourceFile.String())
		if err != nil {
			oper.errLog.Add(err, "unable to Stat", oper.relToSource(sourceFile))
			continue
		}
		if !stat.Mode().IsRegular() {
			oper.errLog.Add(err, "file is not a regular file", oper.relToSource(sourceFile))
			continue
		}

		dirInfo.DiskUsage += stat.Size()
	}

	for _, ch := range dirInfo.Children.Array() {
		dirInfo.DiskUsage += ch.DiskUsage
	}
}

const (
	_        = iota // ignore first value by assigning to blank identifier
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
)

func DirSizeExpr(size int64) string {
	var pref string
	var chunk int64
	if size >= GB/4 {
		pref = "Gb"
		chunk = GB
	} else if size >= MB/4 {
		pref = "Mb"
		chunk = MB
	} else {
		pref = "Kb"
		chunk = KB
	}
	amt := float64(size) / float64(chunk)
	return fmt.Sprintf("%5.1f %v", amt, pref)
}

func (oper Names) processDeleteFlag(path Path) bool {
	result := oper.deleteFlag
	if result {
		if path.IsDir() {
			path.DeleteDirectoryM("~$")
		} else {
			path.DeleteFileM()
		}
	}
	return result
}

func (oper Names) GetHelp(bp *BasePrinter) {
	bp.Pr("Examine filenames; source <source dir> [clean_log]")
}

func (oper Names) examineFilename(p Path) {
	oper.deleteFlag = false
	base := p.Base()

	// See https://en.wikipedia.org/wiki/Tilde
	if WindowsTempPattern.MatchString(base) {
		switch oper.config.Microsoft() {
		default:
			Die("unsupported option:", oper.config.Microsoft())
		case Ignore:
			break
		case Warn:
			oper.errLog.Add(Warning, "Word:", oper.relToSource(p))
		case Delete:
			oper.errLog.Add(Warning, "Deleting Word:", oper.relToSource(p))
			oper.deleteFlag = true
		}
		return
	}

	if !oper.pattern.MatchString(base) {
		summary := oper.highlightStrangeCharacters(base)
		oper.errLog.Add(Warning, "Chars:", summary, "in", oper.relToSource(p))
	}
}

func (oper Names) highlightStrangeCharacters(str string) string {
	// I was doing a binary search, but I found out that due to utf-8, some chars (runes)
	// are different lengths; so just build up the substring from the left until we find the problem
	sb := strings.Builder{}
	sbPost := strings.Builder{}

	problemFound := false
	prob := ""
	for _, ch := range str {
		if !problemFound {
			sb.WriteRune(ch)
			prob = sb.String()
			if !oper.pattern.MatchString(prob) {
				problemFound = true
			}
		} else {
			sbPost.WriteRune(ch)
		}
	}
	return Quoted(sb.String() + "<<<" + sbPost.String())
}

func (oper Names) printDiskUsage(dirInfo DirInfo) {
	Pr(Spaces(dirInfo.Depth*4), DirSizeExpr(dirInfo.DiskUsage), ":", dirInfo.Path.Base())
	if dirInfo.Depth >= int(oper.config.Depth()) {
		return
	}
	for _, ch := range dirInfo.Children.Array() {
		oper.printDiskUsage(ch)
	}
}
