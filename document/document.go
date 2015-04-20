package document

type Document struct {
    Values  map[string]interface{} `json:"values"`
    gloabId uint64
}
