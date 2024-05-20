package gen

import (
	. "github.com/jpsember/golang-base/base"
)

type sNamesConfig struct {
	source          string
	cleanLog        bool
	log             string
	pattern         string
	maxProblems     int
	microsoft       MicrosoftOpt
	verboseProblems bool
	depth           int
}

type NamesConfigBuilderObj struct {
	// We embed the static struct
	sNamesConfig
}

type NamesConfigBuilder = *NamesConfigBuilderObj

// ---------------------------------------------------------------------------------------
// NamesConfig interface
// ---------------------------------------------------------------------------------------

type NamesConfig interface {
	DataClass
	Source() string
	CleanLog() bool
	Log() string
	Pattern() string
	MaxProblems() int
	Microsoft() MicrosoftOpt
	VerboseProblems() bool
	Depth() int
	Build() NamesConfig
	ToBuilder() NamesConfigBuilder
}

var DefaultNamesConfig = newNamesConfig()

// Convenience method to get a fresh builder.
func NewNamesConfig() NamesConfigBuilder {
	return DefaultNamesConfig.ToBuilder()
}

// Construct a new static object, with fields initialized appropriately
func newNamesConfig() NamesConfig {
	var m = sNamesConfig{}
	m.log = "problems.txt"
	m.pattern = "^(\\w|-|\\.|\\x20|,|\\(|\\)|\\+|&|\\$|:|'|#|\\[|\\]|=|!|;|@)+$"
	m.microsoft = DefaultMicrosoftOpt
	m.depth = 255
	return &m
}

// ---------------------------------------------------------------------------------------
// Implementation of static (built) object
// ---------------------------------------------------------------------------------------

func (v *sNamesConfig) Source() string {
	return v.source
}

func (v *sNamesConfig) CleanLog() bool {
	return v.cleanLog
}

func (v *sNamesConfig) Log() string {
	return v.log
}

func (v *sNamesConfig) Pattern() string {
	return v.pattern
}

func (v *sNamesConfig) MaxProblems() int {
	return v.maxProblems
}

func (v *sNamesConfig) Microsoft() MicrosoftOpt {
	return v.microsoft
}

func (v *sNamesConfig) VerboseProblems() bool {
	return v.verboseProblems
}

func (v *sNamesConfig) Depth() int {
	return v.depth
}

func (v *sNamesConfig) Build() NamesConfig {
	// This is already the immutable (built) version.
	return v
}

func (v *sNamesConfig) ToBuilder() NamesConfigBuilder {
	return &NamesConfigBuilderObj{sNamesConfig: *v}
}

func (v *sNamesConfig) ToJson() JSEntity {
	var m = NewJSMap()
	m.Put(NamesConfig_Source, v.source)
	m.Put(NamesConfig_CleanLog, v.cleanLog)
	m.Put(NamesConfig_Log, v.log)
	m.Put(NamesConfig_Pattern, v.pattern)
	m.Put(NamesConfig_MaxProblems, v.maxProblems)
	m.Put(NamesConfig_Microsoft, v.microsoft.String())
	m.Put(NamesConfig_VerboseProblems, v.verboseProblems)
	m.Put(NamesConfig_Depth, v.depth)
	return m
}

func (v *sNamesConfig) Parse(source JSEntity) DataClass {
	var s = source.AsJSMap()
	var n = newNamesConfig().(*sNamesConfig)
	n.source = s.OptString(NamesConfig_Source, "")
	n.cleanLog = s.OptBool(NamesConfig_CleanLog, false)
	n.log = s.OptString(NamesConfig_Log, "problems.txt")
	n.pattern = s.OptString(NamesConfig_Pattern, "^(\\w|-|\\.|\\x20|,|\\(|\\)|\\+|&|\\$|:|'|#|\\[|\\]|=|!|;|@)+$")
	n.maxProblems = s.OptInt(NamesConfig_MaxProblems, 0)
	n.microsoft = DefaultMicrosoftOpt.ParseFrom(s, NamesConfig_Microsoft)
	n.verboseProblems = s.OptBool(NamesConfig_VerboseProblems, false)
	n.depth = s.OptInt(NamesConfig_Depth, 255)
	return n
}

func (v *sNamesConfig) String() string {
	var x = v.ToJson().AsJSMap()
	return PrintJSEntity(x, true)
}

// ---------------------------------------------------------------------------------------
// Implementation of builder
// ---------------------------------------------------------------------------------------

func (v NamesConfigBuilder) Source() string {
	return v.source
}

func (v NamesConfigBuilder) CleanLog() bool {
	return v.cleanLog
}

func (v NamesConfigBuilder) Log() string {
	return v.log
}

func (v NamesConfigBuilder) Pattern() string {
	return v.pattern
}

func (v NamesConfigBuilder) MaxProblems() int {
	return v.maxProblems
}

func (v NamesConfigBuilder) Microsoft() MicrosoftOpt {
	return v.microsoft
}

func (v NamesConfigBuilder) VerboseProblems() bool {
	return v.verboseProblems
}

func (v NamesConfigBuilder) Depth() int {
	return v.depth
}

func (v NamesConfigBuilder) SetSource(source string) NamesConfigBuilder {
	v.source = source
	return v
}

func (v NamesConfigBuilder) SetCleanLog(cleanLog bool) NamesConfigBuilder {
	v.cleanLog = cleanLog
	return v
}

func (v NamesConfigBuilder) SetLog(log string) NamesConfigBuilder {
	v.log = log
	return v
}

func (v NamesConfigBuilder) SetPattern(pattern string) NamesConfigBuilder {
	v.pattern = pattern
	return v
}

func (v NamesConfigBuilder) SetMaxProblems(maxProblems int) NamesConfigBuilder {
	v.maxProblems = maxProblems
	return v
}

func (v NamesConfigBuilder) SetMicrosoft(microsoft MicrosoftOpt) NamesConfigBuilder {
	v.microsoft = microsoft
	return v
}

func (v NamesConfigBuilder) SetVerboseProblems(verboseProblems bool) NamesConfigBuilder {
	v.verboseProblems = verboseProblems
	return v
}

func (v NamesConfigBuilder) SetDepth(depth int) NamesConfigBuilder {
	v.depth = depth
	return v
}

func (v NamesConfigBuilder) Build() NamesConfig {
	// Construct a copy of the embedded static struct
	var b = v.sNamesConfig
	return &b
}

func (v NamesConfigBuilder) ToBuilder() NamesConfigBuilder {
	return v
}

func (v NamesConfigBuilder) ToJson() JSEntity {
	return v.Build().ToJson()
}

func (v NamesConfigBuilder) Parse(source JSEntity) DataClass {
	return DefaultNamesConfig.Parse(source)
}

func (v NamesConfigBuilder) String() string {
	return v.Build().String()
}

const NamesConfig_Source = "source"
const NamesConfig_CleanLog = "clean_log"
const NamesConfig_Log = "log"
const NamesConfig_Pattern = "pattern"
const NamesConfig_MaxProblems = "max_problems"
const NamesConfig_Microsoft = "microsoft"
const NamesConfig_VerboseProblems = "verbose_problems"
const NamesConfig_Depth = "depth"

// Convenience method to parse a NamesConfig from a JSMap
func ParseNamesConfig(jsmap JSEntity) NamesConfig {
	m := jsmap.(JSMap)
	return DefaultNamesConfig.Parse(m).(NamesConfig)
}
