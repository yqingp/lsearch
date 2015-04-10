package store

type MmtrieNode struct {
    key     uint8
    nchilds uint8
    data    int
    childs  int
}

func (self *MmtrieNode) setKey(k byte) {
    self.key = k
    self.nchilds = 0
    self.childs = 0
    self.data = 0
}

func (self *MmtrieNode) nodeCopy(old MmtrieNode) {
    self.childs = old.childs
    self.data = old.data
    self.key = old.key
    self.nchilds = old.nchilds
}
