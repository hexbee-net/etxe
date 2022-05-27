package etx

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"
)

var (
	needsOctalPrefix = regexp.MustCompile(`^0([\d_]+)$`)
	heredocDelimiter = regexp.MustCompile(`^<<([-]?)(\w+)\n$`)
)

// Value is a scalar, list or map.
type Value struct {
	ASTNode

	Comments []string     `parser:"(@Comment [ NewLine ])*" json:"comments,omitempty"`
	Null     bool         `parser:"(  @'null'"               json:"null,omitempty"`
	Bool     *ValueBool   `parser:" | @('true' | 'false')"   json:"bool,omitempty"`
	Number   *ValueNumber `parser:" | @Number"               json:"number,omitempty"`
	Str      *ValueString `parser:" | @@"                    json:"str,omitempty"`
	Ident    *Ident       `parser:" | @@"                    json:"ident,omitempty"`
	Heredoc  *Heredoc     `parser:" | @@"                    json:"heredoc,omitempty"`
	List     *ValueList   `parser:" | @@"                    json:"list,omitempty"`
	Map      *ValueMap    `parser:" | @@ )"                  json:"map,omitempty"`
}

func (v *Value) Clone() *Value {
	if v == nil {
		return nil
	}

	out := Value{
		ASTNode:  v.ASTNode.Clone(),
		Comments: cloneStrings(v.Comments),
		Null:     v.Null,
		Bool:     v.Bool.Clone(),
		Number:   v.Number.Clone(),
		Str:      v.Str.Clone(),
		Ident:    v.Ident.Clone(),
		Heredoc:  v.Heredoc.Clone(),
		List:     v.List.Clone(),
		Map:      v.Map.Clone(),
	}

	if v.Heredoc != nil {
		s := *v.Heredoc
		out.Heredoc = &s
	}

	return &out
}

func (v *Value) Children() (children []Node) {
	if v.Bool != nil {
		children = append(children, v.Bool)
	}

	if v.Number != nil {
		children = append(children, v.Number)
	}

	if v.Str != nil {
		children = append(children, v.Str)
	}

	if v.Ident != nil {
		children = append(children, v.Ident)
	}

	if v.Heredoc != nil {
		children = append(children, v.Heredoc)
	}

	if v.List != nil {
		children = append(children, v.List)
	}

	if v.Map != nil {
		children = append(children, v.Map)
	}

	return
}

func (v Value) String() string {
	switch {
	case v.Null == true:
		return "null"

	case v.Bool != nil:
		return fmt.Sprintf("%v", *v.Bool)

	case v.Number != nil:
		return v.Number.String()

	case v.Str != nil:
		return v.Str.String()

	case v.Ident != nil:
		return v.Ident.String()

	case v.Heredoc != nil:
		return v.Heredoc.String()

	case v.List != nil:
		return v.List.String()

	case v.Map != nil:
		return v.Map.String()

	default:
		panic(repr.String(v, repr.Hide(lexer.Position{})))
	}
}

// /////////////////////////////////////

// ValueBool represents a parsed boolean value.
type ValueBool bool

func (v *ValueBool) Capture(values []string) error {
	*v = values[0] == "true"

	return nil
}

func (v *ValueBool) Clone() *ValueBool {
	if v == nil {
		return nil
	}

	out := *v

	return &out
}

func (v *ValueBool) Children() (children []Node) {
	return
}

// /////////////////////////////////////

// ValueNumber of arbitrary precision.
type ValueNumber struct {
	Value  *big.Float
	Source string
}

// Capture override because big.Float doesn't directly support
// 0-prefix octal parsing ¯\_(ツ)_/¯.
func (v *ValueNumber) Capture(values []string) error {
	v.Source = values[0]

	sm := needsOctalPrefix.FindStringSubmatch(v.Source)
	if sm != nil {
		v.Source = "0o" + sm[1]
	}

	v.Value = big.NewFloat(0)
	if _, _, err := v.Value.Parse(v.Source, 0); err != nil {
		return fmt.Errorf("failed to parse number value '%s': %w", v.Source, err)
	}

	return nil
}

func (v *ValueNumber) Clone() *ValueNumber {
	if v == nil {
		return nil
	}

	out := &ValueNumber{
		Source: v.Source,
	}

	if v.Value != nil {
		out.Value = big.NewFloat(0)
		out.Value.Copy(v.Value)
	}

	return out
}

func (v *ValueNumber) Children() (children []Node) {
	return
}

func (v ValueNumber) String() string {
	if v.Source == "" {
		return v.Value.String()
	}

	return v.Source
}

// /////////////////////////////////////

type Heredoc struct {
	ASTNode

	Delimiter HeredocDelimiter   `parser:"@Heredoc"       json:"delimiter,omitempty"`
	Fragments []*HeredocFragment `parser:"@@* HeredocEnd" json:"fragments,omitempty"`
}

func (v *Heredoc) Clone() *Heredoc {
	if v == nil {
		return nil
	}

	return &Heredoc{
		ASTNode:   v.ASTNode.Clone(),
		Delimiter: *v.Delimiter.Clone(),
		Fragments: cloneCollection(v.Fragments),
	}
}

func (v *Heredoc) Children() (children []Node) {
	for _, item := range v.Fragments {
		children = append(children, item)
	}

	return
}

func (v Heredoc) String() string {
	var sb strings.Builder

	mustFprintf(&sb, "<<%v", v.Delimiter.String())

	for _, fragment := range v.Fragments {
		sb.WriteString(fragment.String())
	}

	sb.WriteString("\n")
	sb.WriteString(v.Delimiter.Delimiter)

	return sb.String()
}

type HeredocDelimiter struct {
	LeadingTabs bool   `json:"leading_tabs"`
	Delimiter   string `json:"delimiter"`
}

func (v *HeredocDelimiter) Capture(values []string) error {
	sm := heredocDelimiter.FindStringSubmatch(values[0])
	if sm == nil {
		panic("missing heredoc delimiter")
	}

	if sm[1] != "" {
		v.LeadingTabs = true
	}

	v.Delimiter = sm[2]

	return nil
}

func (v *HeredocDelimiter) Clone() *HeredocDelimiter {
	return &HeredocDelimiter{
		LeadingTabs: v.LeadingTabs,
		Delimiter:   v.Delimiter,
	}
}

func (v HeredocDelimiter) String() string {
	if v.Delimiter == "" {
		panic("empty heredoc delimiter")
	}

	if v.LeadingTabs {
		return fmt.Sprintf("-%s", v.Delimiter)
	}

	return v.Delimiter
}

type HeredocFragment struct {
	ASTNode

	Expr      *Expr  `parser:"(  Expr @@ ExprEnd       " json:"expr,omitempty"`
	Directive *Expr  `parser:" | Directive @@ ExprEnd  " json:"directive,omitempty"`
	Text      string `parser:" | @(Body|EOL|NonExpr)+ )" json:"text,omitempty"`
}

func (f *HeredocFragment) Clone() *HeredocFragment {
	if f == nil {
		return nil
	}

	return &HeredocFragment{
		ASTNode:   f.ASTNode.Clone(),
		Expr:      f.Expr.Clone(),
		Directive: f.Directive.Clone(),
		Text:      f.Text,
	}
}

func (f *HeredocFragment) Children() (children []Node) {
	if f.Expr != nil {
		children = append(children, f.Expr)
	}

	if f.Directive != nil {
		children = append(children, f.Directive)
	}

	return
}

func (f HeredocFragment) String() string {
	switch {
	case f.Expr != nil:
		return fmt.Sprintf("${ %s }", f.Expr)
	case f.Directive != nil:
		return fmt.Sprintf("%%{ %s }", f.Directive)
	case f.Text != "":
		return f.Text
	default:
		return ""
	}
}

// /////////////////////////////////////

type ValueString struct {
	ASTNode

	Fragment []*StringFragment `parser:"String @@* StringEnd" json:"fragment,omitempty"`
}

func (v *ValueString) Clone() *ValueString {
	if v == nil {
		return nil
	}

	return &ValueString{
		ASTNode:  v.ASTNode.Clone(),
		Fragment: cloneCollection(v.Fragment),
	}
}

func (v *ValueString) Children() (children []Node) {
	for _, item := range v.Fragment {
		children = append(children, item)
	}

	return
}

func (v ValueString) String() string {
	var sb strings.Builder

	sb.WriteString(`"`)

	for _, f := range v.Fragment {
		sb.WriteString(f.String())
	}

	sb.WriteString(`"`)

	return sb.String()
}

type StringFragment struct {
	ASTNode

	Escaped   string `parser:"(  @Escaped"                           json:"escaped,omitempty"`
	Unicode   string `parser:" | Unicode@(UnicodeLong|UnicodeShort)" json:"unicode,omitempty"`
	Expr      *Expr  `parser:" | Expr @@ ExprEnd"                    json:"expr,omitempty"`
	Directive *Expr  `parser:" | Directive @@ ExprEnd"               json:"directive,omitempty"`
	Text      string `parser:" | @(Char|Quote|NonExpr))"             json:"text,omitempty"`
}

func (f *StringFragment) Clone() *StringFragment {
	if f == nil {
		return nil
	}

	return &StringFragment{
		ASTNode:   f.ASTNode.Clone(),
		Escaped:   f.Escaped,
		Unicode:   f.Unicode,
		Expr:      f.Expr.Clone(),
		Directive: f.Directive.Clone(),
		Text:      f.Text,
	}
}

func (f *StringFragment) Children() (children []Node) {
	if f.Expr != nil {
		children = append(children, f.Expr)
	}

	if f.Directive != nil {
		children = append(children, f.Directive)
	}

	return
}

func (f StringFragment) String() string {
	switch {
	case f.Escaped != "":
		return fmt.Sprintf("\\%s", f.Escaped)
	case f.Unicode != "":
		return fmt.Sprintf("\\u%s", f.Unicode)
	case f.Expr != nil:
		return fmt.Sprintf("${%s}", f.Expr)
	case f.Directive != nil:
		return fmt.Sprintf("%%{%s}", f.Directive)
	case f.Text != "":
		return f.Text
	default:
		return ""
	}
}

// /////////////////////////////////////

type ValueList struct {
	ASTNode

	Items []*Value `parser:"'[' NewLine? [ @@ ( NewLine? ',' NewLine? @@? )* ','? ] NewLine? ']'"   json:"items,omitempty"`
}

func (v *ValueList) Clone() *ValueList {
	if v == nil {
		return nil
	}

	return &ValueList{
		ASTNode: v.ASTNode.Clone(),
		Items:   cloneCollection(v.Items),
	}
}

func (v *ValueList) Children() (children []Node) {
	for _, item := range v.Items {
		children = append(children, item)
	}

	return
}

func (v ValueList) String() string {
	if len(v.Items) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")

	for _, e := range v.Items {
		sb.WriteString(indent(e.String(), indentationChar))
		sb.WriteString(",\n")
	}

	sb.WriteString("]")

	return sb.String()
}

// /////////////////////////////////////

type ValueMap struct {
	ASTNode

	Entries []*MapEntry `parser:"'{' NewLine? [ @@ ( NewLine? ',' NewLine? @@? )* ','? ] NewLine? '}'" json:"entries,omitempty"`
}

func (v *ValueMap) Clone() *ValueMap {
	if v == nil {
		return nil
	}

	return &ValueMap{
		ASTNode: v.ASTNode.Clone(),
		Entries: cloneCollection(v.Entries),
	}
}

func (v *ValueMap) Children() (children []Node) {
	for _, item := range v.Entries {
		children = append(children, item)
	}

	return
}

func (v ValueMap) String() string {
	if len(v.Entries) == 0 {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")

	for _, e := range v.Entries {
		sb.WriteString(indent(fmt.Sprintf("%s: %s", e.Key, e.Value), indentationChar))
		sb.WriteString(",\n")
	}

	sb.WriteString("}")

	return sb.String()
}

// MapEntry represents a key+value in a map.
type MapEntry struct {
	ASTNode

	Comments []string `parser:"(@Comment [ NewLine ])*" json:"comments,omitempty"`
	Key      Value    `parser:"@@ ':'"    json:"key"`
	Value    Value    `parser:"@@"        json:"value"`
}

func (v *MapEntry) Clone() *MapEntry {
	if v == nil {
		return nil
	}

	return &MapEntry{
		ASTNode:  v.ASTNode.Clone(),
		Comments: cloneStrings(v.Comments),
		Key:      *v.Key.Clone(),
		Value:    *v.Value.Clone(),
	}
}

func (v *MapEntry) Children() (children []Node) {
	children = append(children, &v.Key)
	children = append(children, &v.Value)

	return
}

func (v MapEntry) String() string {
	return fmt.Sprintf("%v: %v", v.Key, v.Value)
}
