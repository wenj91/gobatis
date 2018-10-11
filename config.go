package gobatis

import (
	"log"
	"sync"
)

type config struct {
	mappedStmts map[string]*node
	mu          sync.Mutex
}

func (this *config) put(id string, n *node) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	if _, ok := this.mappedStmts[id]; ok {
		return false
	}

	this.mappedStmts[id] = n
	return true
}

func (this *config) getMappedStmt(id string) *mappedStmt {
	rootNode, ok := this.mappedStmts[id]
	if !ok {
		log.Fatalln("Can not find id:", id, "mapped stmt")
	}

	resultType := ""
	resultTypeAttr, ok := rootNode.Attrs["resultType"]
	if ok {
		resultType = resultTypeAttr.Value
	}

	sn := createSqlNode(rootNode.Elements...)

	ds := &dynamicSqlSource{
		sqlNode: sn,
	}

	return &mappedStmt{
		sqlSource:  ds,
		resultType: resultType,
	}
}
