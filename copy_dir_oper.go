package main

import (
	. "github.com/jpsember/golang-base/app"
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/files"
	. "golang-diskutils/gen"
	"io"
	"os"
	"path/filepath"
)

type CopyDirOper struct {
	BaseObject
	errLog     ErrLog
	sourcePath Path
	destPath   Path
	errCount   int
	config     CopyDirConfig
	verifyTs   int
}

func (oper *CopyDirOper) GetArguments() DataClass {
	return DefaultCopyDirConfig
}

func (oper *CopyDirOper) ArgsFileMustExist() bool {
	return false
}

func (oper *CopyDirOper) AcceptArguments(a DataClass) {
	oper.config = a.(CopyDirConfig)
}

func (oper *CopyDirOper) UserCommand() string {
	return "copydir"
}

func procPath(desc string, expr string) (Path, string) {
	var err error
	problem := ""
	result := EmptyPath
	for {
		if expr == "" {
			problem = "path is empty"
			break
		}
		absPath, err := filepath.Abs(expr)
		if err != nil {
			break
		}
		result, err = NewPath(absPath)
		if err != nil {
			break
		}
		break
	}
	if err != nil {
		problem = err.Error()
	}
	if problem != "" {
		problem = desc + "; problem: " + problem
	}
	return result, problem
}

func (oper *CopyDirOper) Perform(app *App) {
	oper.SetVerbose(app.Verbose())

	{
		var operSourceDir, operDestDir Path
		problem := ""
		for {
			operSourceDir, problem = procPath("Source directory", oper.config.Source())
			if problem == "" {
				operDestDir, problem = procPath("Target directory", oper.config.Dest())
			}
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
		oper.destPath = operDestDir
	}

	oper.errLog = NewErrLog(oper.config.Log())
	oper.errLog.Clean = oper.config.CleanLog()

	dirStack := NewArray[Path]()
	depthStack := NewArray[int]()
	dirStack.Add(oper.sourcePath)
	depthStack.Add(0)

	sourcePrefixLen := len(oper.sourcePath.String())
	targetPrefix := oper.destPath.String()

	for dirStack.NonEmpty() {
		dir := dirStack.Pop()
		depth := depthStack.Pop()

		// Make target directory if it doesn't already exist
		targetDir := NewPathM(targetPrefix + dir.String()[sourcePrefixLen:])
		err := targetDir.MkDirs()
		if err != nil {
			oper.errLog.Add(err, "unable to make directory", dir)
			continue
		}
		dirEntries, err := os.ReadDir(dir.String())
		if err != nil {
			oper.errLog.Add(err, "unable to read directory contents", dir)
			continue
		}

		for _, dirEntry := range dirEntries {
			nm := dirEntry.Name()

			sourceFile := dir.JoinM(nm)

			// Check if source is a symlink.  If so, skip it.
			srcFileInfo, err := os.Lstat(sourceFile.String())
			if err != nil {
				oper.errLog.Add(err, "unable to get Lstat for", sourceFile)
				continue
			}
			if srcFileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				continue
			}

			if windowsTempPattern.MatchString(sourceFile.Base()) {
				if !oper.config.RetainMicrosoft() {
					oper.errLog.Add(Warning, "skipping Word backup file", sourceFile)
					continue
				}
			}

			sourceFileSuffix := sourceFile.String()[sourcePrefixLen:]
			targetFile := NewPathM(targetPrefix + sourceFileSuffix)

			targetFileExists := targetFile.Exists()
			// If target file already exists, verify it is the same type (dir or file) as source
			if targetFileExists {
				if sourceFile.IsDir() != targetFile.IsDir() {
					oper.errLog.Add(err, "source is not same file/dir type as target:", sourceFile, INDENT,
						"vs", targetFile)
					continue
				}

			}

			if sourceFile.IsDir() {
				dirStack.Add(sourceFile)
				depthStack.Add(depth + 1)
				continue
			}

			sourceFileStat, err := os.Stat(sourceFile.String())
			if err != nil {
				oper.errLog.Add(err, "getting Stat", sourceFile)
				continue
			}
			if !sourceFileStat.Mode().IsRegular() {
				oper.errLog.Add(err, "source file is not a regular file", sourceFile)
				continue
			}

			const timestampApproxEqualMs = 1500

			modifiedTime := sourceFileStat.ModTime()

			action := "copying"

			if targetFileExists {
				// Only continue if source is newer
				targetFileStat, err := os.Stat(targetFile.String())
				if err != nil {
					oper.errLog.Add(err, "getting Stat", targetFile)
					continue
				}

				sourceEpochMs := modifiedTime.UnixMilli()
				targetTime := targetFileStat.ModTime()
				targetEpochMs := targetTime.UnixMilli()
				// There might be a slight roundoff error with the timestamps
				if targetEpochMs+timestampApproxEqualMs >= sourceEpochMs {
					continue
				}
				action = "updating"
			}

			oper.Log(DepthDots(depth, action, sourceFileSuffix))

			err = copyFileContents(sourceFile, targetFile)
			if err != nil {
				oper.errLog.Add(err, action, "file contents", sourceFile, targetFile)
				continue
			}

			err = os.Chtimes(targetFile.String(), modifiedTime, modifiedTime)
			if err != nil {
				oper.errLog.Add(err, "unable to set modified time", targetFile)
				continue
			}

			if oper.verifyTs < 10 {
				oper.verifyTs++
				targetFileStat, err := os.Stat(targetFile.String())
				CheckOk(err, targetFile)
				newTargetTime := targetFileStat.ModTime()
				diff := newTargetTime.UnixMilli() - modifiedTime.UnixMilli()
				CheckArg(Abs(diff) <= timestampApproxEqualMs, "target time still different")
			}
		}
	}
	oper.errLog.PrintSummary()
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Copies file.  If destination exists, its contents will be replaced.
func copyFileContents(srcp, dstp Path) (err error) {
	src := srcp.String()
	dst := dstp.String()

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (oper *CopyDirOper) GetHelp(bp *BasePrinter) {
	bp.Pr("Copy a directory; source <source dir> dest <dest dir> [clean_log]")
}

func addCopyDirOper(app *App) {
	var oper = &CopyDirOper{}
	oper.ProvideName(oper)
	app.RegisterOper(AssertJsonOper(oper))
}
