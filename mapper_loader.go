package gobatis

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func lookupMapper(paths ...string) []string {
	var fs []string

	for _, p := range paths {
		_ = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".xml") {
				p, err = filepath.Abs(path)
				if nil != err {
					return err
				}
				fs = append(fs, p)
				return nil
			}
			return nil
		})
	}
	return fs
}

func loadingMapper(paths ...string) *mapper {
	fs := lookupMapper(paths...)

	mp := newMapper()

	for _, f := range fs {
		r, e := os.Open(f)
		if nil != e {
			continue
		}

		rootNode := parse(r)

		if rootNode.Name != "mapper" {
			log.Fatalln("mapper xml must start with `mapper` tag, please check your xml mapper!")
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
						log.Fatalln("No id for:", childNode.Name, "Id must be not null, please check your xml mapper!")
					}

					fid := namespace + childNode.Id
					if ok := mp.put(fid, &childNode); !ok {
						log.Fatalln("repeat id for:", fid, "Please check your xml mapper!")
					}

				case "sql":
					if childNode.Id == "" {
						log.Fatalln("no id for:", childNode.Name, "Id must be not null, please check your xml mapper!")
					}

					fid := namespace + childNode.Id
					if ok := mp.put(fid, &childNode); !ok {
						log.Fatalln("repeat id for:", fid, "Please check your xml mapper!")
					}
				}
			}
		}
		_ = r.Close()
	}

	return mp
}
