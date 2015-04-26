package analyzer

import (
	"testing"
)

func TestA(t *testing.T) {
	an := NewAnalyzer()
	an.Init()
	words := an.Analyze("测试微博索引")
	for k, _ := range words {
		t.Log(k)
	}
}
