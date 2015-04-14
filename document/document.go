package document

type Document struct {
    Id      int64
    fields  map[string]interface{}
    IndexId int
}
