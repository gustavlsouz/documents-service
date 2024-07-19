package common

import (
	"database/sql"
	"encoding/json"
)

type JsonNullBool struct {
	sql.NullBool
}

func (value JsonNullBool) MarshalJSON() ([]byte, error) {
	if value.Valid {
		return json.Marshal(value.Bool)
	} else {
		return json.Marshal(nil)
	}
}

func (value *JsonNullBool) UnmarshalJSON(data []byte) error {
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		value.Valid = true
		value.Bool = *x
	} else {
		value.Valid = false
	}
	return nil
}
