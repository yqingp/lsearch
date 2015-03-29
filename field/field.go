package field

import (
	"github.com/yqingp/lsearch/analyzer"
)

type Filed struct {
	Id             int
	Name           string
	CreatedAt      int64
	FieldType      FieldType
	SearchAnalyzer *analyzer.Analyzer
	IndexAnalyzer  *analyzer.Analyzer
	IsIndex        bool
}
