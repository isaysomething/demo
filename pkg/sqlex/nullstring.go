package sqlex

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func ToNullString(s string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	} else {
		return json.Marshal(nil)
	}
}

func (n NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		n.Valid = true
		n.String = *s
	} else {
		n.Valid = false
	}
	return nil
}
