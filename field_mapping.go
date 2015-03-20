package lsearch

type FieldMapping struct {
	Name     string
	Type     FieldType
	Analyzer string
	IsStore  bool
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
