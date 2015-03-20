package lsearch

import (
// "bytes"
// "encoding/binary"
// "errors"
// "fmt"
// "log"
// "os"
// "time"
)

// type IndexField struct {
// 	Id          uint16
// 	Name        string
// 	CreatedTime int64
// 	FieldType   uint8
// 	IsIndex     bool
// 	IsStore     bool
// 	IsAnalyze   bool
// }

type Index struct {
	Id   uint16
	Name string
	// Fields           map[uint16]IndexField
	// FiledsName       map[string]IndexField
	ShardNum         uint16
	ReplicaNum       uint16
	TotalTokenNum    uint32
	TotalDocumentNum uint64
	TotalFieldNum    uint16
	// Shards           map[int]Shard
	// Replicas         map[int]Replica
}

// type Indexes struct {
// 	TotalIndexNum uint16
// 	indexes       map[string]*Index
// 	file          *os.File
// 	fname         string
// }

// var IndexesMap Indexes

// var indexesMetaMagicNum int64 = 0x1234567123

// func Init() {
// 	IndexesMap = Indexes{TotalIndexNum: 0, indexes: make(map[string]*Index)}
// 	f, err := os.Create("indexes.meta")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	IndexesMap.file = f
// 	IndexesMap.fname = "indexes.meta"
// }

// func NewIndex(name string, fields map[string]interface{}) (index *Index, err error) {
// 	if index = IndexesMap.indexes[name]; index != nil {
// 		err = errors.New(ERROR_INDEX_NAME_HAS_EXIST)
// 		return
// 	}

// 	index = &Index{Name: name,
// 		Id: IndexesMap.TotalIndexNum + 1,
// 		// ShardNum:   DeafultShardNum,
// 		// ReplicaNum: DefaultReplicaNum,
// 	}

// 	initIndexFields(index, fields)
// 	fmt.Println("%s", index)
// 	IndexesMap.TotalIndexNum++
// 	IndexesMap.indexes[index.Name] = index
// 	updateIndexesMata()
// 	return
// }

// func initIndexFields(index *Index, fields map[string]interface{}) (err error) {
// 	var currentTotalFieldNum uint16 = 1
// 	for name, v := range fields {
// 		oneField := IndexField{
// 			Name:        name,
// 			Id:          currentTotalFieldNum,
// 			CreatedTime: time.Now().Unix(),
// 		}
// 		if fieldType, ok := v.(uint8); !ok {
// 			log.Fatal(ERROR_TYPE_ERROR)
// 		} else {
// 			oneField.FieldType = fieldType
// 		}

// 		index.Fields = make(map[uint16]IndexField)
// 		index.Fields[oneField.Id] = oneField
// 		index.FiledsName = make(map[string]IndexField)
// 		index.FiledsName[oneField.Name] = oneField
// 		index.TotalFieldNum++
// 		currentTotalFieldNum++
// 	}
// 	return
// }

// func updateIndexesMata() (int, error) {
// 	log.Println(IndexesMap.TotalIndexNum)
// 	if IndexesMap.TotalIndexNum == 0 {
// 		log.Println("total index Num is 0")
// 		return 0, nil
// 	}

// 	var err error
// 	buf := &bytes.Buffer{}

// 	err = binary.Write(buf, binary.LittleEndian, indexesMetaMagicNum)
// 	if err != nil {
// 		log.Fatal(err)
// 		return 0, err
// 	}

// 	err = binary.Write(buf, binary.LittleEndian, (IndexesMap.TotalIndexNum))
// 	if err != nil {
// 		log.Fatal(err)
// 		return 0, err
// 	}

// 	for _, v := range IndexesMap.indexes {
// 		err = binary.Write(buf, binary.LittleEndian, v.Id)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}

// 		err = binary.Write(buf, binary.LittleEndian, v.ShardNum)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}
// 		err = binary.Write(buf, binary.LittleEndian, v.ReplicaNum)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}
// 		err = binary.Write(buf, binary.LittleEndian, v.TotalDocumentNum)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}
// 		err = binary.Write(buf, binary.LittleEndian, v.TotalTokenNum)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}
// 		err = binary.Write(buf, binary.LittleEndian, v.TotalFieldNum)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}

// 		err = binary.Write(buf, binary.LittleEndian, uint16(len(v.Name)))
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}

// 		_, err = buf.WriteString(v.Name)
// 		if err != nil {
// 			log.Fatal(err)
// 			return 0, err
// 		}
// 		logger.Println("11")

// 		for _, vv := range v.Fields {
// 			err = binary.Write(buf, binary.LittleEndian, vv.Id)
// 			if err != nil {
// 				log.Fatal(err)
// 				return 0, err
// 			}

// 			err = binary.Write(buf, binary.LittleEndian, uint16(len(vv.Name)))
// 			if err != nil {
// 				log.Fatal(err)
// 				return 0, err
// 			}

// 			_, err = buf.WriteString(vv.Name)
// 			if err != nil {
// 				log.Fatal(err)
// 				return 0, err
// 			}

// 			err = binary.Write(buf, binary.LittleEndian, vv.CreatedTime)
// 			if err != nil {
// 				log.Fatal(err)
// 				return 0, err
// 			}

// 			err = binary.Write(buf, binary.LittleEndian, vv.FieldType)
// 			if err != nil {
// 				log.Fatal(err)
// 				return 0, err
// 			}

// 			if vv.IsAnalyze {
// 				err = binary.Write(buf, binary.LittleEndian, uint8(1))
// 				if err != nil {
// 					log.Fatal(err)
// 					return 0, err
// 				}
// 			} else {
// 				err = binary.Write(buf, binary.LittleEndian, uint8(0))
// 				if err != nil {
// 					log.Fatal(err)
// 					return 0, err
// 				}
// 			}

// 			// err = binary.Write(buf, binary.LittleEndian, vv.IsIndex)
// 			// if err != nil {
// 			// 	return 0, err
// 			// }

// 			// err = binary.Write(buf, binary.LittleEndian, vv.IsStore)
// 			// if err != nil {
// 			// 	return 0, err
// 			// }
// 		}

// 		IndexesMap.file.Write(buf.Bytes())
// 	}
// 	logger.Println("11s")
// 	// log.Println(math.MaxUint8)
// 	// log.Println(math.MaxUint16)
// 	// log.Println(math.MaxUint32)
// 	// log.Println(math.MaxUint64)

// 	return 0, nil

// }
