package gen

import (
	. "github.com/jpsember/golang-base/base"
)

// See https://www.sohamkamani.com/golang/enums/

type MicrosoftOpt int

const (
	Ignore MicrosoftOpt = iota
	Warn
	Delete
)

const DefaultMicrosoftOpt = Ignore

var MicrosoftOptEnumInfo = NewEnumInfo("ignore warn delete")

func (x MicrosoftOpt) String() string {
	return MicrosoftOptEnumInfo.EnumNames[x]
}

func (x MicrosoftOpt) ParseFrom(m JSMap, key string) MicrosoftOpt {
	return MicrosoftOpt(ParseEnumFromMap(MicrosoftOptEnumInfo, m, key, int(DefaultMicrosoftOpt)))
}
