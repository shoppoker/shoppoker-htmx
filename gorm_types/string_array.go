package gorm_types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	srcTyped, ok := src.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to StringArray", src, src)
	}
	return json.Unmarshal([]byte(srcTyped), a)
}

func (a StringArray) Value() (driver.Value, error) {
	s, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return string(s), nil
}
