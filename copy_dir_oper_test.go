package main

import (
	. "github.com/jpsember/golang-base/base"
	. "github.com/jpsember/golang-base/json"
	"github.com/jpsember/golang-base/jt"
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
	j.GenerateSubdirs("source", jsmap)
	jm := jt.DirSummary(j.GetTestResultsDir())
	Pr(jm.CompactString())
	j.AssertGenerated()
}
