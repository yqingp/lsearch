package engine

import (
    "encoding/json"
    "errors"
    "github.com/yqingp/lsearch/index"
    "github.com/yqingp/lsearch/query"
    "time"
)

/*
{
	:name => "weibo",
	:query => {
		:text => "weibo",
		:from => 1,
		:limit => n
	}
}
*/

type SearchRequest struct {
    Name         string `json:"name"`
    RequestStart time.Time
    Status       chan bool
    Results      interface{}
    Duration     string
    Index        *index.Index
    Query        *query.Query `json:"query,omitempty"`
    Error        error
}

func (e *Engine) Search(body []byte) interface{} {
    searchRequest, err := ParseSearchRequest(body)

    response := struct {
        Took    string      `json:"took"`
        Results interface{} `json:"results"`
        Error   string      `json:"error,omitempty"`
    }{}

    if err != nil {
        response.Results = ""
        response.Took = "0ms"
        response.Error = err.Error()

        return response
    }

    searchRequest.RequestStart = time.Now()
    searchRequest.Status = make(chan bool)

    if err = searchRequest.Valid(); err != nil {
        return err
    }

    index, ok := e.indexes[searchRequest.Name]

    if !ok {
        response.Results = ""
        response.Took = "0ms"
        response.Error = "Index Not Found"
    }
    searchRequest.Index = index

    // Logger.Println(indexRequest)

    e.SearchRequests <- searchRequest
    <-searchRequest.Status
    searchRequest.Duration = time.Now().Sub(searchRequest.RequestStart).String()

    response.Results = searchRequest.Results
    response.Took = searchRequest.Duration

    return response
}

func ParseSearchRequest(body []byte) (*SearchRequest, error) {
    request := &SearchRequest{}
    if err := json.Unmarshal(body, request); err != nil {
        return nil, errors.New("decode request error")
    }

    return request, nil
}

func (i *SearchRequest) Valid() error {
    if i.Name == "" || i.Query == nil || i.Query.Text == "" {
        return errors.New("Search Request Error")
    }

    i.Duration = time.Now().Sub(i.RequestStart).String()

    return nil
}
