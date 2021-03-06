package etx

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
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

	Null    bool         `parser:"(  @'null'"             json:"null,omitempty"`
	Bool    *ValueBool   `parser:" | @('true' | 'false')" json:"bool,omitempty"`
	Number  *ValueNumber `parser:" | @Number"             json:"number,omitempty"`
	Str     *ValueString `parser:" | @@"                  json:"str,omitempty"`
	Heredoc *Heredoc     `parser:" | @@"                  json:"heredoc,omitempty"`
	List    *ValueList   `parser:" | @@"                  json:"list,omitempty"`
	Map     *ValueMap    `parser:" | @@ )"                json:"map,omitempty"`
}

func (v *Value) Clone() *Value {
	if v == nil {
		return nil
	}

	out := Value{
		ASTNode: v.ASTNode.Clone(),
		Null:    v.Null,
		Bool:    v.Bool.Clone(),
		Number:  v.Number.Clone(),
		Str:     v.Str.Clone(),
		Heredoc: v.Heredoc.Clone(),
		List:    v.List.Clone(),
		Map:     v.Map.Clone(),
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

func (v Value) FormattedString() string {
	var sb strings.Builder

	switch {
	case v.Null == true:
		sb.WriteString("null")

	case v.Bool != nil:
		mustFprintf(&sb, "%s", v.Bool.FormattedString())

	case v.Number != nil:
		sb.WriteString(v.Number.FormattedString())

	case v.Str != nil:
		sb.WriteString(v.Str.FormattedString())

	case v.Heredoc != nil:
		sb.WriteString(v.Heredoc.FormattedString())

	case v.List != nil:
		sb.WriteString(v.List.FormattedString())

	case v.Map != nil:
		sb.WriteString(v.Map.FormattedString())

	default:
		panic(repr.String(v, repr.Hide(lexer.Position{})))
	}

	return sb.String()
}

// /////////////////////////////////////

// ValueBool represents a parsed boolean value.
type ValueBool struct {
	ASTNode

	Value bool
}

func (v *ValueBool) Capture(values []string) error {
	v.Value = values[0] == "true"

	return nil
}

func (v *ValueBool) Clone() *ValueBool {
	if v == nil {
		return nil
	}

	out := &ValueBool{
		ASTNode: v.ASTNode.Clone(),
		Value:   v.Value,
	}

	return out
}

func (v *ValueBool) Children() (children []Node) {
	return
}

func (v ValueBool) FormattedString() string {
	return strconv.FormatBool(v.Value)
}

// /////////////////////////////////////

// ValueNumber of arbitrary precision.
type ValueNumber struct {
	ASTNode

	Value  *big.Float
	Source string
}

// Capture override because big.Float doesn't directly support
// 0-prefix octal parsing ??\_(???)_/??.
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
		ASTNode: v.ASTNode.Clone(),
		Source:  v.Source,
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

func (v ValueNumber) FormattedString() string {
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

func (v Heredoc) FormattedString() string {
	var sb strings.Builder

	mustFprintf(&sb, "<<%s", v.Delimiter.FormattedString())

	for _, fragment := range v.Fragments {
		sb.WriteString(fragment.FormattedString())
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

func (v HeredocDelimiter) FormattedString() string {
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

func (f HeredocFragment) FormattedString() string {
	switch {
	case f.Expr != nil:
		return fmt.Sprintf("${ %s }", f.Expr.FormattedString())
	case f.Directive != nil:
		return fmt.Sprintf("%%{ %s }", f.Directive.FormattedString())
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

func (v ValueString) FormattedString() string {
	var sb strings.Builder

	sb.WriteString(`"`)

	for _, f := range v.Fragment {
		sb.WriteString(f.FormattedString())
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

func (f StringFragment) FormattedString() string {
	switch {
	case f.Escaped != "":
		return fmt.Sprintf("\\%s", f.Escaped)
	case f.Unicode != "":
		return fmt.Sprintf("\\u%s", f.Unicode)
	case f.Expr != nil:
		return fmt.Sprintf("${%s}", f.Expr.FormattedString())
	case f.Directive != nil:
		return fmt.Sprintf("%%{%s}", f.Directive.FormattedString())
	case f.Text != "":
		return f.Text
	default:
		return ""
	}
}

// /////////////////////////////////////

type ValueList struct {
	ASTNode

	Items []*ListItem `parser:"'[' [ LF+ ] @@*  ']'"   json:"items,omitempty"`
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

func (v ValueList) FormattedString() string {
	if len(v.Items) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")

	for _, e := range v.Items {
		sb.WriteString(indent(e.FormattedString(), indentationChar))
	}

	sb.WriteString("]")

	return sb.String()
}

type ListItem struct {
	ASTNode

	EmptyLine string   `parser:"(   @LF+             " json:"empty_line,omitempty"`
	Value     *Expr    `parser:"  | ( @@ ','? LF? )  " json:"value,omitempty"`
	Comment   *Comment `parser:"  | @@              )" json:"comment,omitempty"`
}

func (v *ListItem) Clone() *ListItem {
	if v == nil {
		return nil
	}

	return &ListItem{
		ASTNode:   v.ASTNode.Clone(),
		EmptyLine: v.EmptyLine,
		Value:     v.Value.Clone(),
		Comment:   v.Comment.Clone(),
	}
}

func (v *ListItem) Children() (children []Node) {
	if v.Value != nil {
		children = append(children, v.Value)
	}

	if v.Comment != nil {
		children = append(children, v.Comment)
	}

	return
}

func (v ListItem) FormattedString() string {
	switch {
	case v.EmptyLine != "":
		return v.EmptyLine
	case v.Comment != nil:
		return v.Comment.FormattedString()
	case v.Value != nil:
		return fmt.Sprintf("%s,\n", v.Value.FormattedString())
	default:
		panic("item not set")
	}
}

// /////////////////////////////////////

type ValueMap struct {
	ASTNode

	Items []*MapItem `parser:"'{' [ LF+ ] @@* '}'" json:"items,omitempty"`
}

func (v *ValueMap) Clone() *ValueMap {
	if v == nil {
		return nil
	}

	return &ValueMap{
		ASTNode: v.ASTNode.Clone(),
		Items:   cloneCollection(v.Items),
	}
}

func (v *ValueMap) Children() (children []Node) {
	for _, item := range v.Items {
		children = append(children, item)
	}

	return
}

func (v ValueMap) FormattedString() string {
	if len(v.Items) == 0 {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")

	for _, e := range v.Items {
		sb.WriteString(indent(fmt.Sprintf("%s: %s", e.Key.FormattedString(), e.Value.FormattedString()), indentationChar))
		sb.WriteString(",\n")
	}

	sb.WriteString("}")

	return sb.String()
}

// MapItem represents a key+value in a map.
type MapItem struct {
	ASTNode

	EmptyLine string   `parser:"(   @LF+             " json:"empty_line,omitempty"`
	Key       *MapKey  `parser:"  | ( @@ '='         " json:"key"`
	Value     *Expr    `parser:"      @@ ','? LF? )  " json:"value"`
	Comment   *Comment `parser:"  | @@              )" json:"comment,omitempty"`
}

func (v *MapItem) Clone() *MapItem {
	if v == nil {
		return nil
	}

	return &MapItem{
		ASTNode:   v.ASTNode.Clone(),
		EmptyLine: v.EmptyLine,
		Key:       v.Key.Clone(),
		Value:     v.Value.Clone(),
		Comment:   v.Comment.Clone(),
	}
}

func (v *MapItem) Children() (children []Node) {
	if v.Key != nil {
		children = append(children, v.Key)
	}

	if v.Value != nil {
		children = append(children, v.Value)
	}

	if v.Comment != nil {
		children = append(children, v.Comment)
	}

	return
}

func (v MapItem) FormattedString() string {
	switch {
	case v.EmptyLine != "":
		return v.EmptyLine
	case v.Comment != nil:
		return v.Comment.FormattedString()
	case v.Key != nil && v.Value != nil:
		return fmt.Sprintf("%s = %s", v.Key.FormattedString(), v.Value.FormattedString())
	default:
		panic("item not set")
	}
}

// MapKey represent a key in a MapItem.
type MapKey struct {
	ASTNode

	Ident *Ident       `parser:"(   @@  " json:"ident,omitempty"`
	Str   *ValueString `parser:"  | @@ )" json:"str,omitempty"`
}

func (v *MapKey) Clone() *MapKey {
	if v == nil {
		return nil
	}

	return &MapKey{
		ASTNode: v.ASTNode.Clone(),
		Ident:   v.Ident.Clone(),
		Str:     v.Str.Clone(),
	}
}

func (v *MapKey) Children() (children []Node) {
	switch {
	case v.Ident != nil:
		children = append(children, v.Ident)
	case v.Str != nil:
		children = append(children, v.Str)
	}

	return
}

func (v MapKey) FormattedString() string {
	switch {
	case v.Ident != nil:
		return v.Ident.FormattedString()
	case v.Str != nil:
		return v.Str.FormattedString()
	default:
		panic("key is not set")
	}
}
