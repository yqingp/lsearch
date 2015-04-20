package field

import ()

type Filed struct {
    Id        int       `json:"id"`
    CreatedAt int64     `json:"created_at"`
    Name      string    `json:"name"`
    FieldType FieldType `json:"type"`
    IsIndex   bool      `json:"is_index,omitempty"`
}

func (f Filed) Valid() (isValid bool) {
    if f.Name == "" {
        return false
    }

    if f.FieldType != 0 && f.FieldType != 1 {
        return false
    }

    return true
}
