package document

import (
    "encoding/json"
    "github.com/yqingp/lsearch/analyzer"
    "github.com/yqingp/lsearch/field"
)

// {
// 	id => "1", values => {a => "1", b => 2}
// }

type Document struct {
    gloabId uint64
    id      string `json:"id"`
    tokens  map[string]string
    Values  map[string]interface{} `json:"values"`
}

func (d *Document) SetGloabId(id uint64) {
    d.gloabId = id
}

func (d *Document) GloabId() uint64 {
    return d.gloabId
}

func (d *Document) Id() string {
    return d.id
}

func (d Document) Tokens() map[string]string {

    return d.Tokens()
}

func Validate(documents []Document, fields []field.Filed) bool {

    newFields := map[string]field.Filed{}

    for _, v := range fields {
        newFields[v.Name] = v
    }

    for _, doc := range documents {
        if doc.Id() == "" {
            return false
        }
        for k, v := range doc.Values {
            field, ok := newFields[k]
            if !ok {
                return false
            }

            if !field.IsValidValue(v) {
                return false
            }
        }
    }

    return true
}

func (d *Document) Analyze(analyzer *analyzer.Analyzer) {
    for _, v := range d.Values {
        value, ok := v.(string)
        if ok && value != "" {
            words := analyzer.Analyze(value)
            for k, _ := range words {
                d.tokens[k] = ""
            }
        }
    }
}

func (d Document) Encode() ([]byte, error) {
    data, err := json.Marshal(d)
    if err != nil {
        return nil, err
    }

    return data, nil
}

func (d *Document) Decode(data []byte) error {

    err := json.Unmarshal(data, d)
    if err != nil {
        return err
    }

    return nil
}
