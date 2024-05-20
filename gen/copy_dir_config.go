package gen

import (
	. "github.com/jpsember/golang-base/base"
)

type sCopyDirConfig struct {
	source          string
	dest            string
	cleanLog        bool
	log             string
	retainMicrosoft bool
}

type CopyDirConfigBuilderObj struct {
	// We embed the static struct
	sCopyDirConfig
}

type CopyDirConfigBuilder = *CopyDirConfigBuilderObj

// ---------------------------------------------------------------------------------------
// CopyDirConfig interface
// ---------------------------------------------------------------------------------------

type CopyDirConfig interface {
	DataClass
	Source() string
	Dest() string
	CleanLog() bool
	Log() string
	RetainMicrosoft() bool
	Build() CopyDirConfig
	ToBuilder() CopyDirConfigBuilder
}

var DefaultCopyDirConfig = newCopyDirConfig()

// Convenience method to get a fresh builder.
func NewCopyDirConfig() CopyDirConfigBuilder {
	return DefaultCopyDirConfig.ToBuilder()
}

// Construct a new static object, with fields initialized appropriately
func newCopyDirConfig() CopyDirConfig {
	var m = sCopyDirConfig{}
	m.log = "problems.txt"
	return &m
}

// ---------------------------------------------------------------------------------------
// Implementation of static (built) object
// ---------------------------------------------------------------------------------------

func (v *sCopyDirConfig) Source() string {
	return v.source
}

func (v *sCopyDirConfig) Dest() string {
	return v.dest
}

func (v *sCopyDirConfig) CleanLog() bool {
	return v.cleanLog
}

func (v *sCopyDirConfig) Log() string {
	return v.log
}

func (v *sCopyDirConfig) RetainMicrosoft() bool {
	return v.retainMicrosoft
}

func (v *sCopyDirConfig) Build() CopyDirConfig {
	// This is already the immutable (built) version.
	return v
}

func (v *sCopyDirConfig) ToBuilder() CopyDirConfigBuilder {
	return &CopyDirConfigBuilderObj{sCopyDirConfig: *v}
}

func (v *sCopyDirConfig) ToJson() JSEntity {
	var m = NewJSMap()
	m.Put(CopyDirConfig_Source, v.source)
	m.Put(CopyDirConfig_Dest, v.dest)
	m.Put(CopyDirConfig_CleanLog, v.cleanLog)
	m.Put(CopyDirConfig_Log, v.log)
	m.Put(CopyDirConfig_RetainMicrosoft, v.retainMicrosoft)
	return m
}

func (v *sCopyDirConfig) Parse(source JSEntity) DataClass {
	var s = source.AsJSMap()
	var n = newCopyDirConfig().(*sCopyDirConfig)
	n.source = s.OptString(CopyDirConfig_Source, "")
	n.dest = s.OptString(CopyDirConfig_Dest, "")
	n.cleanLog = s.OptBool(CopyDirConfig_CleanLog, false)
	n.log = s.OptString(CopyDirConfig_Log, "problems.txt")
	n.retainMicrosoft = s.OptBool(CopyDirConfig_RetainMicrosoft, false)
	return n
}

func (v *sCopyDirConfig) String() string {
	var x = v.ToJson().AsJSMap()
	return PrintJSEntity(x, true)
}

// ---------------------------------------------------------------------------------------
// Implementation of builder
// ---------------------------------------------------------------------------------------

func (v CopyDirConfigBuilder) Source() string {
	return v.source
}

func (v CopyDirConfigBuilder) Dest() string {
	return v.dest
}

func (v CopyDirConfigBuilder) CleanLog() bool {
	return v.cleanLog
}

func (v CopyDirConfigBuilder) Log() string {
	return v.log
}

func (v CopyDirConfigBuilder) RetainMicrosoft() bool {
	return v.retainMicrosoft
}

func (v CopyDirConfigBuilder) SetSource(source string) CopyDirConfigBuilder {
	v.source = source
	return v
}

func (v CopyDirConfigBuilder) SetDest(dest string) CopyDirConfigBuilder {
	v.dest = dest
	return v
}

func (v CopyDirConfigBuilder) SetCleanLog(cleanLog bool) CopyDirConfigBuilder {
	v.cleanLog = cleanLog
	return v
}

func (v CopyDirConfigBuilder) SetLog(log string) CopyDirConfigBuilder {
	v.log = log
	return v
}

func (v CopyDirConfigBuilder) SetRetainMicrosoft(retainMicrosoft bool) CopyDirConfigBuilder {
	v.retainMicrosoft = retainMicrosoft
	return v
}

func (v CopyDirConfigBuilder) Build() CopyDirConfig {
	// Construct a copy of the embedded static struct
	var b = v.sCopyDirConfig
	return &b
}

func (v CopyDirConfigBuilder) ToBuilder() CopyDirConfigBuilder {
	return v
}

func (v CopyDirConfigBuilder) ToJson() JSEntity {
	return v.Build().ToJson()
}

func (v CopyDirConfigBuilder) Parse(source JSEntity) DataClass {
	return DefaultCopyDirConfig.Parse(source)
}

func (v CopyDirConfigBuilder) String() string {
	return v.Build().String()
}

const CopyDirConfig_Source = "source"
const CopyDirConfig_Dest = "dest"
const CopyDirConfig_CleanLog = "clean_log"
const CopyDirConfig_Log = "log"
const CopyDirConfig_RetainMicrosoft = "retain_microsoft"

// Convenience method to parse a CopyDirConfig from a JSMap
func ParseCopyDirConfig(jsmap JSEntity) CopyDirConfig {
	m := jsmap.(JSMap)
	return DefaultCopyDirConfig.Parse(m).(CopyDirConfig)
}
