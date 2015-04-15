package field

import ()

type Filed struct {
    Id        int
    Name      string `json:"name"`
    CreatedAt int64
    FieldType FieldType `json:"type"`
    IsIndex   bool      `json:"is_index",omitempty`
}
