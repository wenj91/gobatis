package gobatis

import (
	"io"
	"strings"

	"gopkg.in/yaml.v2"
)

func buildMapperConfig(r io.Reader) *mapperConfig {
	rootNode := parse(r)

	conf := &mapperConfig{
		mappedStmts: make(map[string]*node),
		mappedSql:   make(map[string]*node),
		cache:       make(map[string]*mappedStmt),
	}

	if rootNode.Name != "mapper" {
		LOG.Error("Mapper xml must start with `mapper` tag, please check your xml mapperConfig!")
		panic("Mapper xml must start with `mapper` tag, please check your xml mapperConfig!")
	}

	namespace := ""
	if val, ok := rootNode.Attrs["namespace"]; ok {
		nStr := strings.TrimSpace(val.Value)
		if nStr != "" {
			nStr += "."
		}
		namespace = nStr
	}

	for _, elem := range rootNode.Elements {
		if elem.ElementType == eleTpNode {
			childNode := elem.Val.(node)
			switch childNode.Name {
			case "select", "update", "insert", "delete":
				if childNode.Id == "" {
					LOG.Error("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
					panic("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
				}

				fid := namespace + childNode.Id
				if ok := conf.put(fid, &childNode); !ok {
					LOG.Error("Repeat id for:" + fid + "Please check your xml mapperConfig!")
					panic("Repeat id for:" + fid + "Please check your xml mapperConfig!")
				}

			case "sql":
				if childNode.Id == "" {
					LOG.Error("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
					panic("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
				}

				fid := namespace + childNode.Id
				if ok := conf.putSql(fid, &childNode); !ok {
					LOG.Error("Repeat id for:" + fid + "Please check your xml mapperConfig!")
					panic("Repeat id for:" + fid + "Please check your xml mapperConfig!")
				}
			}
		}
	}

	return conf
}

func buildDbConfig(ymlStr string) *DBConfig {
	dbconf := &DBConfig{}
	err := yaml.Unmarshal([]byte(ymlStr), &dbconf)
	if err != nil {
		panic("error: " + err.Error())
	}

	return dbconf
}
