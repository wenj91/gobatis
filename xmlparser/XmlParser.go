package xmlparser

import (
	"encoding/xml"
	"os"
	"github.com/wenj91/gobatis/stack"
	"github.com/wenj91/gobatis/tools/tagutil"
	"log"
)

const (
	ELE_TP_STRING = 0
	ELE_TP_NODE   = 1
)

type node struct {
	Name     string
	Attr     []xml.Attr
	Elements []element
}

type element struct {
	ElementType int
	StrVal      string
	NodeVal     node
}

func Parse(file *os.File) []node {
	parser := xml.NewDecoder(file)
	nodes := make([]node, 0)

	st := stack.NewStack()
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			if (tagutil.IsContains(name)) {
				attr := elmt.Attr
				nd := node{
					Name:     name,
					Attr:     attr,
					Elements: make([]element, 0),
				}
				st.Push(nd)
			}
		case xml.EndElement:
			if st.Len() > 0 {
				n := st.Pop().(node)
				if st.Len()>0 {
					e := element{
						ElementType: ELE_TP_NODE,
						NodeVal: n,
					}
					pn := st.Pop().(node)
					els := pn.Elements
					els = append(els, e)
					pn.Elements = els
					st.Push(pn)
				}else{
					nodes = append(nodes, n)
				}
			}
		case xml.CharData:
			if st.Len()>0 {
				n := st.Pop().(node)
				if (tagutil.IsContains(n.Name)) {
					bytes := xml.CharData(t)
					e := element{
						ElementType:ELE_TP_STRING,
						StrVal: string([]byte(bytes)),
					}
					els := n.Elements
					els = append(els, e)

					n.Elements = els
				}
				st.Push(n)
			}

		case xml.ProcInst:
			log.Println("xml:ProcInst")
		case xml.Directive:
			log.Println("xml:Directive")
		case xml.Comment:
			log.Println("xml:Comment")
		default:
			log.Println("xml:Unknown")
		}
	}

	//
	//bs , _ := json.Marshal(nodes)
	//fmt.Println("node:", len(nodes), string(bs))
	return nodes
}


