package analyzer

import (
	"github.com/wangbin/jiebago"
)


// Analyzer struct
type Analyzer struct {
	Name string
}


// NewAnalyzer 创建
func NewAnalyzer() *Analyzer {
	return &Analyzer{Name: "jieba"}
}


// Init  初始化分词器
func (a *Analyzer) Init() {
	jiebago.SetDictionary("")
}


// Analyze  分词
func (a *Analyzer) Analyze(text string) map[string]bool {
	words := make(map[string]bool)
	ch := jiebago.CutForSearch(text, true)
	for word := range ch {
		if len(word) <= 1 {
			continue
		}
		words[word] = true
	}

	return words
}
