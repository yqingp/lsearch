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
