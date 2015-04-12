package analyzer

import (
    "testing"
)

func TestA(t *testing.T) {
    an := NewAnalyzer()
    an.Init()
    words := an.Analyze("测试分词")
    for _, v := range words {
        t.Log(v)
    }
}
