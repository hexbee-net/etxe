package etx

import (
	"fmt"
	"strings"

	"github.com/alecthomas/repr"
)

type Type struct {
	ASTNode

	Label  string      `parser:"'type' @Ident"                        json:"label"`
	Enum   *TypeEnum   `parser:"(   Enum   '{' [ NewLine+ ] @@ "      json:"enum,omitempty"`
	Object *TypeObject `parser:"  | Object '{' [ NewLine+ ] @@ ) '}'" json:"object,omitempty"`
}

func (n *Type) Clone() *Type {
	if n == nil {
		return nil
	}

	out := &Type{
		ASTNode: n.ASTNode.Clone(),
		Label:   n.Label,
		Enum:    n.Enum.Clone(),
		Object:  n.Object.Clone(),
	}

	return out
}

func (n *Type) Children() (children []Node) {
	if n.Enum != nil {
		children = append(children, n.Enum)
	}

	if n.Object != nil {
		children = append(children, n.Object)
	}

	return
}

func (n Type) String() string {
	var sb strings.Builder

	mustFprintf(&sb, "type %s ", n.Label)

	switch {
	case n.Enum != nil:
		mustFprintf(&sb, "enum {%v\n}", indent(n.Enum.String(), indentationChar))
	case n.Object != nil:
		mustFprintf(&sb, "object {%v\n}", indent(n.Object.String(), indentationChar))
	default:
		panic(repr.String(n, repr.Hide(Position{})))
	}

	return sb.String()
}

// /////////////////////////////////////

type TypeEnum struct {
	ASTNode

	Items []*TypeEnumItem `parser:"@@*"  json:"items,omitempty"`
}

func (n *TypeEnum) Clone() *TypeEnum {
	if n == nil {
		return nil
	}

	return &TypeEnum{
		ASTNode: n.ASTNode.Clone(),
		Items:   cloneCollection(n.Items),
	}
}

func (n *TypeEnum) Children() (children []Node) {
	for _, item := range n.Items {
		children = append(children, item)
	}

	return
}

func (n TypeEnum) String() string {
	var sb strings.Builder

	maxLabelLength := 0
	for _, item := range n.Items {
		if l := len(item.Label); l > maxLabelLength {
			maxLabelLength = l
		}
	}

	for _, item := range n.Items {
		fillLength := maxLabelLength - len(item.Label)
		mustFprintf(&sb, "\n%v:%v %v", item.Label, strings.Repeat(" ", fillLength), item.Value)
	}

	return sb.String()
}

// /////////////////////////////////////

type TypeEnumItem struct {
	ASTNode

	EmptyLine string   `parser:"(   @NewLine+         " json:"empty_line,omitempty"`
	Comment   *Comment `parser:"  | @@                " json:"comment,omitempty"`
	Label     string   `parser:"  | ( @Ident ':'      " json:"label"`
	Value     Expr     `parser:"      @@ [NewLine] ) )" json:"value"`
}

func (n *TypeEnumItem) Clone() *TypeEnumItem {
	if n == nil {
		return nil
	}

	return &TypeEnumItem{
		ASTNode:   n.ASTNode.Clone(),
		Comment:   n.Comment.Clone(),
		Label:     n.Label,
		Value:     *n.Value.Clone(),
		EmptyLine: n.EmptyLine,
	}
}

func (n *TypeEnumItem) Children() (children []Node) {
	if n.Comment != nil {
		children = append(children, n.Comment)
	}

	children = append(children, &n.Value)

	return
}

func (n *TypeEnumItem) String() string {
	switch {
	case n.Comment != nil:
		return n.Comment.String()
	case n.EmptyLine != "":
		return n.EmptyLine
	default:
		return fmt.Sprintf("%v: %v", n.Label, n.Value.String())
	}
}

// /////////////////////////////////////

type TypeObject struct {
	ASTNode

	Items []*TypeObjectItem `parser:"@@*" json:"items,omitempty"`
}

func (n *TypeObject) Clone() *TypeObject {
	if n == nil {
		return nil
	}

	out := &TypeObject{
		ASTNode: n.ASTNode.Clone(),
		Items:   cloneCollection(n.Items),
	}

	return out
}

func (n *TypeObject) Children() (children []Node) {
	for _, item := range n.Items {
		children = append(children, item)
	}

	return
}

func (n TypeObject) String() string {
	var sb strings.Builder

	maxLabelLength := 0
	for _, item := range n.Items {
		if l := len(item.Label); l > maxLabelLength {
			maxLabelLength = l
		}
	}

	for _, item := range n.Items {
		fillLength := maxLabelLength - len(item.Label)
		mustFprintf(&sb, "\n%v:%v %v", item.Label, strings.Repeat(" ", fillLength), item.Type)
	}

	return sb.String()
}

// /////////////////////////////////////

type TypeObjectItem struct {
	ASTNode

	EmptyLine string        `parser:"(   @NewLine+         " json:"empty_line,omitempty"`
	Comment   *Comment      `parser:"  | @@                " json:"comment,omitempty"`
	Label     string        `parser:"  | ( @Ident ':'      " json:"label"`
	Type      ParameterType `parser:"      @@ [NewLine] ) )" json:"type"`
}

func (n *TypeObjectItem) Clone() *TypeObjectItem {
	if n == nil {
		return nil
	}

	return &TypeObjectItem{
		ASTNode:   n.ASTNode.Clone(),
		Comment:   n.Comment.Clone(),
		Label:     n.Label,
		Type:      *n.Type.Clone(),
		EmptyLine: n.EmptyLine,
	}
}

func (n *TypeObjectItem) Children() (children []Node) {
	if n.Comment != nil {
		children = append(children, n.Comment)
	}

	children = append(children, &n.Type)

	return
}

func (n TypeObjectItem) String() string {
	switch {
	case n.Comment != nil:
		return n.Comment.String()
	case n.EmptyLine != "":
		return n.EmptyLine
	default:
		return fmt.Sprintf("%v: %v", n.Label, n.Type.String())
	}
}
