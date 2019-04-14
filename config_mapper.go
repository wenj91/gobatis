package gobatis

import (
	"log"
	"sync"
)

type mapperConfig struct {
	mappedStmts map[string]*node
	mappedSql   map[string]*node
	cache       map[string]*mappedStmt
	mu          sync.Mutex
}

func (this *mapperConfig) put(id string, n *node) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	if _, ok := this.mappedStmts[id]; ok {
		return false
	}

	this.mappedStmts[id] = n
	return true
}

func (this *mapperConfig) putSql(id string, n *node) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	if _, ok := this.mappedSql[id]; ok {
		return false
	}

	this.mappedSql[id] = n
	return true
}

func (this *mapperConfig) getMappedStmt(id string) *mappedStmt {
	if nil == this.cache {
		this.cache = make(map[string]*mappedStmt)
	}
	
	if st, ok := this.cache[id]; ok {
		return st
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	rootNode, ok := this.mappedStmts[id]
	if !ok {
		log.Fatalln("Can not find id:", id, "mapped stmt")
	}

	resultType := ""
	if rootNode.Name == "select" {
		resultTypeAttr, ok := rootNode.Attrs["resultType"]
		if !ok {
			log.Fatalln("Tag `<select>` must have resultType attr!")
		}

		resultType = resultTypeAttr.Value
	}

	sn := createSqlNode(rootNode.Elements...)

	ds := &dynamicSqlSource{}
	ds.sqlNode = sn[0]
	if len(sn) > 1 {
		ds.sqlNode = &mixedSqlNode{
			sqlNodes: sn,
		}
	}

	stmt := &mappedStmt{
		sqlSource:  ds,
		resultType: ResultType(resultType),
	}

	this.cache[id] = stmt

	return stmt
}
