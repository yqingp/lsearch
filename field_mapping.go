package lsearch

type FieldMapping struct {
	Name       string
	Type       string
	Analyzer   string
	IsStore    bool
	IsIndex    bool
	IsAnalyzed bool
}

// func NewDefaultFiledMapping() *FieldMapping {
// 	return &FieldMapping{
// 		Name:       "",
// 		Type:       "",
// 		Analyzer:   "a",
// 		IsStore:    true,
// 		IsIndex:    true,
// 		IsAnalyzed: true,
// 	}
// }
