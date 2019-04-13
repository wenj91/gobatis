package gobatis

import (
	"log"
	"sync"
)

type mapper struct {
	mappedStmts map[string]*node
	mappedSql   map[string]*node
	cache       map[string]*mappedStmt
	mu          sync.Mutex
}

type mappedStmt struct {
	dbType     DbType
	sqlSource  iSqlSource
	resultType ResultType
}

func newMapper() *mapper {
	return &mapper{
		mappedStmts: make(map[string]*node),
		mappedSql:   make(map[string]*node),
		cache:       make(map[string]*mappedStmt),
		mu:          sync.Mutex{},
	}
}

func (m *mapper) put(id string, n *node) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.mappedStmts[id]; ok {
		return false
	}

	m.mappedStmts[id] = n
	return true
}

func (m *mapper) putSql(id string, n *node) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.mappedSql[id]; ok {
		return false
	}

	m.mappedSql[id] = n
	return true
}

func (m *mapper) getMappedStmt(id string) *mappedStmt {
	if st, ok := m.cache[id]; ok {
		return st
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	rootNode, ok := m.mappedStmts[id]
	if !ok {
		log.Fatalln("can not find id:", id, "mapped Executor")
	}

	resultType := ""
	if rootNode.Name == "select" {
		resultTypeAttr, ok := rootNode.Attrs["resultType"]
		if !ok {
			log.Fatalln("tag `<select>` must have resultType attr!")
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

	m.cache[id] = stmt

	return stmt
}
