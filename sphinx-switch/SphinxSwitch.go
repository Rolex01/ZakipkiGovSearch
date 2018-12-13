package sphinxSwitch

import (
	"bitbucket.org/company-one/tender-one/state"
)

const useDublicateStateKey = "sphinxSwitch_useDublicate"

var useDublicate bool = false
var hosts []string

func LoadState() {

	value, err := state.ReadBool(useDublicateStateKey)

	if err == nil {
		useDublicate = value
	}
}

func Switch(toDublicate bool) {
	state.WriteBool(useDublicateStateKey, toDublicate)
	useDublicate = toDublicate
}

func GetHost() string {

	if useDublicate {
		return hosts[1]
	} else {
		return hosts[0]
	}
}

func Init(origHost, dubHost string) {
	hosts = []string{origHost, dubHost}
}
