package helpers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

/*
		@DESC type data to save json on field database
		JSONB Interface for JSONB Field of master product Table
	 |
*/
type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

/*
		@DESC type data to save onlytime on field database
		OnlyTime Interface for OnlyTime Field of master product Table
	 |
*/

// Only time
const TimeFormat = "15:04:05"
const DateFormat = "02-01-2006"

type OnlyTime time.Time

func NewTime(hour, min, sec int) OnlyTime {
	t := time.Date(0, time.January, 1, hour, min, sec, 0, time.UTC)
	return OnlyTime(t)
}

func (t *OnlyTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalText(string(v))
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = OnlyTime(v)
	case nil:
		*t = OnlyTime{}
	default:
		return fmt.Errorf("cannot sql.Scan() MyTime from: %#v", v)
	}
	return nil
}

func (t OnlyTime) Value() (driver.Value, error) {
	return driver.Value(time.Time(t).Format(TimeFormat)), nil
}

func (t *OnlyTime) UnmarshalText(value string) error {
	dd, err := time.Parse(TimeFormat, value)
	if err != nil {
		return err
	}
	*t = OnlyTime(dd)
	return nil
}

func (OnlyTime) GormDataType() string {
	return "TIME"
}
