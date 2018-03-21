package xmlparser

import (
	"github.com/wenj91/gobatis/constants"
	"github.com/wenj91/gobatis/stack"
	"encoding/xml"
	"io"
)

type Node struct {
	Id       string
	Name     string
	Attr     map[string]xml.Attr
	Elements []Element
}

type Element struct {
	ElementType int
	Val         interface{}
}

func Parse(r io.Reader) Node {
	parser := xml.NewDecoder(r)
	var root Node

	st := stack.NewStack()
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement: //tag start
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			attr := elmt.Attr
			attrMap := make(map[string]xml.Attr)
			for _, val := range attr {
				attrMap[val.Name.Local] = val
			}
			node := Node{
				Name:     name,
				Attr:     attrMap,
				Elements: make([]Element, 0),
			}
			for _, val := range attr {
				if val.Name.Local == "id" {
					node.Id = val.Value
				}
			}
			st.Push(node)

		case xml.EndElement: //tag end
			if st.Len() > 0 {
				//cur node
				n := st.Pop().(Node)
				if st.Len() > 0 { //if the root node then append to element
					e := Element{
						ElementType: constants.ELE_TP_NODE,
						Val:         n,
					}

					pn := st.Pop().(Node)
					els := pn.Elements
					els = append(els, e)
					pn.Elements = els
					st.Push(pn)
				} else { //else root = n
					root = n
				}
			}
		case xml.CharData: //tag content
			if st.Len() > 0 {
				n := st.Pop().(Node)

				bytes := xml.CharData(t)
				e := Element{
					ElementType: constants.ELE_TP_STRING,
					Val:         string([]byte(bytes)),
				}
				els := n.Elements
				els = append(els, e)

				n.Elements = els

				st.Push(n)
			}

		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		default:
		}
	}

	return root
}
