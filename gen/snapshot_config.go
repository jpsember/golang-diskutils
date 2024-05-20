package gen

import (
	. "github.com/jpsember/golang-base/base"
)

type sSnapshotConfig struct {
	intervalSeconds    int
	maxImages          int
	outputDir          string
	debugMaxIterations int
	numDevices         int
}

type SnapshotConfigBuilderObj struct {
	// We embed the static struct
	sSnapshotConfig
}

type SnapshotConfigBuilder = *SnapshotConfigBuilderObj

// ---------------------------------------------------------------------------------------
// SnapshotConfig interface
// ---------------------------------------------------------------------------------------

type SnapshotConfig interface {
	DataClass
	IntervalSeconds() int
	MaxImages() int
	OutputDir() string
	DebugMaxIterations() int
	NumDevices() int
	Build() SnapshotConfig
	ToBuilder() SnapshotConfigBuilder
}

var DefaultSnapshotConfig = newSnapshotConfig()

// Convenience method to get a fresh builder.
func NewSnapshotConfig() SnapshotConfigBuilder {
	return DefaultSnapshotConfig.ToBuilder()
}

// Construct a new static object, with fields initialized appropriately
func newSnapshotConfig() SnapshotConfig {
	var m = sSnapshotConfig{}
	m.intervalSeconds = 300
	m.maxImages = 2000
	m.numDevices = 2
	return &m
}

// ---------------------------------------------------------------------------------------
// Implementation of static (built) object
// ---------------------------------------------------------------------------------------

func (v *sSnapshotConfig) IntervalSeconds() int {
	return v.intervalSeconds
}

func (v *sSnapshotConfig) MaxImages() int {
	return v.maxImages
}

func (v *sSnapshotConfig) OutputDir() string {
	return v.outputDir
}

func (v *sSnapshotConfig) DebugMaxIterations() int {
	return v.debugMaxIterations
}

func (v *sSnapshotConfig) NumDevices() int {
	return v.numDevices
}

func (v *sSnapshotConfig) Build() SnapshotConfig {
	// This is already the immutable (built) version.
	return v
}

func (v *sSnapshotConfig) ToBuilder() SnapshotConfigBuilder {
	return &SnapshotConfigBuilderObj{sSnapshotConfig: *v}
}

func (v *sSnapshotConfig) ToJson() JSEntity {
	var m = NewJSMap()
	m.Put(SnapshotConfig_IntervalSeconds, v.intervalSeconds)
	m.Put(SnapshotConfig_MaxImages, v.maxImages)
	m.Put(SnapshotConfig_OutputDir, v.outputDir)
	m.Put(SnapshotConfig_DebugMaxIterations, v.debugMaxIterations)
	m.Put(SnapshotConfig_NumDevices, v.numDevices)
	return m
}

func (v *sSnapshotConfig) Parse(source JSEntity) DataClass {
	var s = source.AsJSMap()
	var n = newSnapshotConfig().(*sSnapshotConfig)
	n.intervalSeconds = s.OptInt(SnapshotConfig_IntervalSeconds, 300)
	n.maxImages = s.OptInt(SnapshotConfig_MaxImages, 2000)
	n.outputDir = s.OptString(SnapshotConfig_OutputDir, "")
	n.debugMaxIterations = s.OptInt(SnapshotConfig_DebugMaxIterations, 0)
	n.numDevices = s.OptInt(SnapshotConfig_NumDevices, 2)
	return n
}

func (v *sSnapshotConfig) String() string {
	var x = v.ToJson().AsJSMap()
	return PrintJSEntity(x, true)
}

// ---------------------------------------------------------------------------------------
// Implementation of builder
// ---------------------------------------------------------------------------------------

func (v SnapshotConfigBuilder) IntervalSeconds() int {
	return v.intervalSeconds
}

func (v SnapshotConfigBuilder) MaxImages() int {
	return v.maxImages
}

func (v SnapshotConfigBuilder) OutputDir() string {
	return v.outputDir
}

func (v SnapshotConfigBuilder) DebugMaxIterations() int {
	return v.debugMaxIterations
}

func (v SnapshotConfigBuilder) NumDevices() int {
	return v.numDevices
}

func (v SnapshotConfigBuilder) SetIntervalSeconds(intervalSeconds int) SnapshotConfigBuilder {
	v.intervalSeconds = intervalSeconds
	return v
}

func (v SnapshotConfigBuilder) SetMaxImages(maxImages int) SnapshotConfigBuilder {
	v.maxImages = maxImages
	return v
}

func (v SnapshotConfigBuilder) SetOutputDir(outputDir string) SnapshotConfigBuilder {
	v.outputDir = outputDir
	return v
}

func (v SnapshotConfigBuilder) SetDebugMaxIterations(debugMaxIterations int) SnapshotConfigBuilder {
	v.debugMaxIterations = debugMaxIterations
	return v
}

func (v SnapshotConfigBuilder) SetNumDevices(numDevices int) SnapshotConfigBuilder {
	v.numDevices = numDevices
	return v
}

func (v SnapshotConfigBuilder) Build() SnapshotConfig {
	// Construct a copy of the embedded static struct
	var b = v.sSnapshotConfig
	return &b
}

func (v SnapshotConfigBuilder) ToBuilder() SnapshotConfigBuilder {
	return v
}

func (v SnapshotConfigBuilder) ToJson() JSEntity {
	return v.Build().ToJson()
}

func (v SnapshotConfigBuilder) Parse(source JSEntity) DataClass {
	return DefaultSnapshotConfig.Parse(source)
}

func (v SnapshotConfigBuilder) String() string {
	return v.Build().String()
}

const SnapshotConfig_IntervalSeconds = "interval_seconds"
const SnapshotConfig_MaxImages = "max_images"
const SnapshotConfig_OutputDir = "output_dir"
const SnapshotConfig_DebugMaxIterations = "debug_max_iterations"
const SnapshotConfig_NumDevices = "num_devices"

// Convenience method to parse a SnapshotConfig from a JSMap
func ParseSnapshotConfig(jsmap JSEntity) SnapshotConfig {
	m := jsmap.(JSMap)
	return DefaultSnapshotConfig.Parse(m).(SnapshotConfig)
}
