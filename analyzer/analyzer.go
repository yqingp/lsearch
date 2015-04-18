package analyzer

import (
	"github.com/wangbin/jiebago"
)

type Analyzer struct {
	Name string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{Name: "jieba"}
}

func (a *Analyzer) Init() {
	jiebago.SetDictionary("")
}

func (a *Analyzer) Analyze(text string) []string {
	var words []string
	ch := jiebago.CutForSearch(text, true)
	for word := range ch {
		words = append(words, word)
	}

	return words
}
