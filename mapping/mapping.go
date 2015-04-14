package lsearch

type IndexMapping struct {
    Name                   string
    DefaultIndexerAnalyzer string
    DeafultSearchAanlyzer  string
    ReplicaNum             int
    ShardNum               int
}

func NewIndexMapping() IndexMapping {
    return IndexMapping{
        Name: "",
        DefaultIndexerAnalyzer: "none",
        DeafultSearchAanlyzer:  "none",
    }
}
