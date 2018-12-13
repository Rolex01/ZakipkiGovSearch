package psqlSwitch

import (
	"database/sql"
	"bitbucket.org/company-one/tender-one/state"
	"bitbucket.org/company-one/tender-one/postgres-driver"
)

const useDublicateStateKey = "psqlHost_useDublicate"

var useDublicate bool = false
var currentDB *sql.DB
var configs []*psql.Database

func LoadState() {

	value, err := state.ReadBool(useDublicateStateKey)

	if err == nil {
		useDublicate = value
	}
}

func Switch(toDublicate bool) error {

	state.WriteBool(useDublicateStateKey, toDublicate)
	useDublicate = toDublicate

	if configs == nil || currentDB == nil {
		return nil
	}

	currentConfigIndex := 0

	if useDublicate {
		currentConfigIndex = 1
	}

	nextDB, err := psql.Init(configs[currentConfigIndex])

	if err != nil {
		return err
	}

	currentDB.Close()
	*currentDB = *nextDB

	return nil
}

func Init(origConfig, dubConfig *psql.Database) (*sql.DB, error) {

	currentConfigs := []*psql.Database{origConfig, dubConfig}
	currentConfigIndex := 0

	if useDublicate {
		currentConfigIndex = 1
	}

	db, err := psql.Init(currentConfigs[currentConfigIndex])

	if err != nil {
		return nil, err
	}

	currentDB = db
	configs = currentConfigs

	return db, nil
}
