package query

type Query struct {
    Text  string `json:"text"`
    From  int    `json:"from"`
    limit int    `json:"limie"`
}
