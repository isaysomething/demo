package sqlex

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func ToNullTime(t time.Time) NullTime {
	return NullTime{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: !t.IsZero(),
		},
	}
}

func (n *NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (n NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		n.Valid = true
		n.Time = *t
	} else {
		n.Valid = false
	}
	return nil
}
