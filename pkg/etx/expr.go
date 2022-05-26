package etx

import (
	"fmt"
	"strings"
)

type Expr struct {
	ASTNode

	Left   *ExprConditional `parser:"(   @@  " json:"left,omitempty"`
	If     *ExprIf          `parser:"  | @@  " json:"if,omitempty"`
	Switch *ExprSwitch      `parser:"  | @@ )" json:"switch,omitempty"`
}

func (e *Expr) Clone() *Expr {
	if e == nil {
		return nil
	}

	return &Expr{
		ASTNode: e.ASTNode.Clone(),
		Left:    e.Left.Clone(),
		If:      e.If.Clone(),
		Switch:  e.Switch.Clone(),
	}
}

func (e *Expr) Children() (children []Node) {
	switch {
	case e.Left != nil:
		children = append(children, e.Left)
	case e.If != nil:
		children = append(children, e.If)
	case e.Switch != nil:
		children = append(children, e.Switch)

	}

	return
}

func (e Expr) String() string {
	switch {
	case e.Left != nil:
		return e.Left.String()
	case e.If != nil:
		return e.If.String()
	case e.Switch != nil:
		return e.Switch.String()
	default:
		panic("expression not set")
	}
}

// /////////////////////////////////////

type ExprIf struct {
	ASTNode

	Condition ExprLogicalOr `parser:"If @@ NewLine? '{' NewLine?"           json:"condition"`
	Left      *Expr         `parser:"@@? NewLine? '}' NewLine?"             json:"left,omitempty"`
	Right     *Expr         `parser:"[ Else '{' NewLine? @@ NewLine? '}' ]" json:"right,omitempty"`
}

func (e *ExprIf) Clone() *ExprIf {
	if e == nil {
		return nil
	}

	return &ExprIf{
		ASTNode:   e.ASTNode.Clone(),
		Condition: *e.Condition.Clone(),
		Left:      e.Left.Clone(),
		Right:     e.Right.Clone(),
	}
}

func (e *ExprIf) Children() (children []Node) {
	children = append(children, &e.Condition)

	if e.Left != nil {
		children = append(children, e.Left)
	}

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprIf) String() string {
	switch {
	case e.Left == nil:
		return fmt.Sprintf("if %s { }", e.Condition)
	case e.Right == nil:
		return fmt.Sprintf("if %s {\n%s\n}", e.Condition, indent(e.Left.String(), indentationChar))
	default:
		return fmt.Sprintf("if %s {\n%s\n} else {\n%s\n}", e.Condition, indent(e.Left.String(), indentationChar), indent(e.Right.String(), indentationChar))

	}
}

// /////////////////////////////////////

type ExprSwitch struct {
	ASTNode

	Selector ExprLogicalOr `parser:"Switch @@ NewLine? '{' NewLine?"   json:"selector"`
	Cases    []*ExprCase   `parser:"(@@ NewLine?)* '}'"                json:"cases,omitempty"`
}

func (e *ExprSwitch) Clone() *ExprSwitch {
	if e == nil {
		return nil
	}

	return &ExprSwitch{
		ASTNode:  e.ASTNode.Clone(),
		Selector: *e.Selector.Clone(),
		Cases:    cloneCollection(e.Cases),
	}
}

func (e *ExprSwitch) Children() (children []Node) {
	children = append(children, &e.Selector)

	for _, item := range e.Cases {
		children = append(children, item)
	}

	return

}

func (e ExprSwitch) String() string {
	switch len(e.Cases) {
	case 0:
		return fmt.Sprintf("switch %s { }", e.Selector)

	default:
		cases := make([]string, 0, len(e.Cases))
		for _, c := range e.Cases {
			cases = append(cases, c.String())
		}

		return fmt.Sprintf("switch %s {\n%s\n}", e.Selector, indent(strings.Join(cases, "\n"), indentationChar))
	}
}

// /////////////////////////////////////

type ExprCase struct {
	ASTNode

	Conditions []*ExprLogicalOr `parser:"(   Case @@ ( ',' @@ )* ':'  "         json:"conditions"`
	Default    bool             `parser:"  | @'default'          ':' )"         json:"default"`
	Expr       *Expr            `parser:"NewLine? '{' NewLine? @@ NewLine? '}'" json:"expr,omitempty"`
}

func (e *ExprCase) Clone() *ExprCase {
	if e == nil {
		return nil
	}

	return &ExprCase{
		ASTNode:    e.ASTNode.Clone(),
		Conditions: cloneCollection(e.Conditions),
		Default:    e.Default,
		Expr:       e.Expr.Clone(),
	}
}

func (e *ExprCase) Children() (children []Node) {
	for _, item := range e.Conditions {
		children = append(children, item)
	}

	if e.Expr != nil {
		children = append(children, e.Expr)
	}

	return
}

func (e ExprCase) String() string {
	switch {
	case e.Conditions != nil:
		conditions := make([]string, 0, len(e.Conditions))
		for _, c := range e.Conditions {
			conditions = append(conditions, c.String())
		}

		return fmt.Sprintf("case %s: {\n%s\n}", strings.Join(conditions, ", "), indent(e.Expr.String(), indentationChar))

	case e.Default:
		return fmt.Sprintf("default: {\n%s\n}", indent(e.Expr.String(), indentationChar))

	default:
		panic("non-default case statement without condition")
	}
}

// /////////////////////////////////////

// ExprConditional is a ternary expression.
//
// Ternaries are bad but necessary for Terraform compatibility, so they are
// included at the expression top level but `if` and `switch` just skip them
// and go straight to the next level.
type ExprConditional struct {
	ASTNode

	Condition   ExprLogicalOr  `parser:"@@"          json:"condition"`
	ConditionOp bool           `parser:"[ @'?'     " json:"condition_op,omitempty"`
	TrueExpr    *ExprLogicalOr `parser:"  @@       " json:"true_expr,omitempty"`
	FalseExpr   *ExprLogicalOr `parser:"  ':' @@  ]" json:"false_expr,omitempty"`
}

func (e *ExprConditional) Clone() *ExprConditional {
	if e == nil {
		return nil
	}

	return &ExprConditional{
		ASTNode:     e.ASTNode.Clone(),
		Condition:   *e.Condition.Clone(),
		ConditionOp: e.ConditionOp,
		TrueExpr:    e.TrueExpr.Clone(),
		FalseExpr:   e.FalseExpr.Clone(),
	}
}

func (e *ExprConditional) Children() (children []Node) {
	children = append(children, &e.Condition)

	if e.TrueExpr != nil {
		children = append(children, e.TrueExpr)
	}

	if e.FalseExpr != nil {
		children = append(children, e.FalseExpr)
	}

	return
}

func (e ExprConditional) String() string {
	if e.ConditionOp {
		switch {
		case e.TrueExpr != nil && e.FalseExpr != nil:
			return fmt.Sprintf("%s ? %s : %s", e.Condition, e.TrueExpr, e.FalseExpr)
		case e.TrueExpr != nil && e.FalseExpr == nil:
			return fmt.Sprintf("%s ? %s : null", e.Condition, e.TrueExpr)
		case e.TrueExpr == nil && e.FalseExpr != nil:
			return fmt.Sprintf("%s ? null : %s", e.Condition, e.FalseExpr)
		default:
			return fmt.Sprintf("%s ? null : null", e.Condition)
		}
	}

	return e.Condition.String()
}

// /////////////////////////////////////

type ExprLogicalOr struct {
	ASTNode

	Left  ExprLogicalAnd `parser:"@@"               json:"left"`
	Op    string         `parser:"[ @OpLogicalOr  " json:"op,omitempty"`
	Right *ExprLogicalOr `parser:"  @@           ]" json:"right,omitempty"`
}

func (e *ExprLogicalOr) Clone() *ExprLogicalOr {
	if e == nil {
		return nil
	}

	return &ExprLogicalOr{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprLogicalOr) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprLogicalOr) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprLogicalAnd struct {
	ASTNode

	Left  ExprBitwiseOr   `parser:"@@"                json:"left"`
	Op    string          `parser:"[ @OpLogicalAnd  " json:"op,omitempty"`
	Right *ExprLogicalAnd `parser:"  @@            ]" json:"right,omitempty"`
}

func (e *ExprLogicalAnd) Clone() *ExprLogicalAnd {
	if e == nil {
		return nil
	}

	return &ExprLogicalAnd{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprLogicalAnd) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprLogicalAnd) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprBitwiseOr struct {
	ASTNode

	Left  ExprBitwiseXor `parser:"@@"               json:"left"`
	Op    string         `parser:"[ @OpBitwiseOr  " json:"op,omitempty"`
	Right *ExprBitwiseOr `parser:"  @@           ]" json:"right,omitempty"`
}

func (e *ExprBitwiseOr) Clone() *ExprBitwiseOr {
	if e == nil {
		return nil
	}

	return &ExprBitwiseOr{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprBitwiseOr) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprBitwiseOr) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprBitwiseXor struct {
	ASTNode

	Left  ExprBitwiseAnd  `parser:"@@"                json:"left"`
	Op    string          `parser:"[ @OpBitwiseXOr  " json:"op,omitempty"`
	Right *ExprBitwiseXor `parser:"  @@            ]" json:"right,omitempty"`
}

func (e *ExprBitwiseXor) Clone() *ExprBitwiseXor {
	if e == nil {
		return nil
	}

	return &ExprBitwiseXor{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprBitwiseXor) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprBitwiseXor) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprBitwiseAnd struct {
	ASTNode

	Left  ExprEquality    `parser:"@@"                json:"left"`
	Op    string          `parser:"[ @OpBitwiseAnd  " json:"op,omitempty"`
	Right *ExprBitwiseAnd `parser:"  @@            ]" json:"right,omitempty"`
}

func (e *ExprBitwiseAnd) Clone() *ExprBitwiseAnd {
	if e == nil {
		return nil
	}

	return &ExprBitwiseAnd{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprBitwiseAnd) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprBitwiseAnd) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprEquality struct {
	ASTNode

	Left  ExprRelational `parser:"@@"                            json:"left"`
	Op    string         `parser:"[ @( OpNotEqual | OpEqual )  " json:"op,omitempty"`
	Right *ExprEquality  `parser:"  @@                        ]" json:"right,omitempty"`
}

func (e *ExprEquality) Clone() *ExprEquality {
	if e == nil {
		return nil
	}

	return &ExprEquality{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprEquality) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprEquality) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprRelational struct {
	ASTNode

	Left  ExprShift       `parser:"@@"                                                       json:"left"`
	Op    string          `parser:"[ @( OpMore | OpMoreOrEqual | OpLess | OpLessOrEqual )  " json:"op,omitempty"`
	Right *ExprRelational `parser:"  @@                                                   ]" json:"right,omitempty"`
}

func (e *ExprRelational) Clone() *ExprRelational {
	if e == nil {
		return nil
	}

	return &ExprRelational{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprRelational) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprRelational) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprShift struct {
	ASTNode

	Left  ExprAdditive `parser:"@@"                                                json:"left"`
	Op    string       `parser:"[ @( OpBitwiseShiftLeft | OpBitwiseShiftRight )  " json:"op,omitempty"`
	Right *ExprShift   `parser:"  @@                                            ]" json:"right,omitempty"`
}

func (e *ExprShift) Clone() *ExprShift {
	if e == nil {
		return nil
	}

	return &ExprShift{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprShift) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprShift) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprAdditive struct {
	ASTNode

	Left  ExprMultiplicative `parser:"@@"                        json:"left"`
	Op    string             `parser:"[ @( OpMinus | OpPlus )  " json:"op,omitempty"`
	Right *ExprAdditive      `parser:"  @@                    ]" json:"right,omitempty"`
}

func (e *ExprAdditive) Clone() *ExprAdditive {
	if e == nil {
		return nil
	}

	return &ExprAdditive{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprAdditive) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprAdditive) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprMultiplicative struct {
	ASTNode

	Left  ExprUnary           `parser:"@@"                                                json:"left,omitempty"`
	Op    string              `parser:"[ @( OpDivision | OpMultiplication | OpModulo )  " json:"op,omitempty"`
	Right *ExprMultiplicative `parser:"  @@                                            ]" json:"right,omitempty"`
}

func (e *ExprMultiplicative) Clone() *ExprMultiplicative {
	if e == nil {
		return nil
	}

	return &ExprMultiplicative{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Op:      e.Op,
		Right:   e.Right.Clone(),
	}
}

func (e *ExprMultiplicative) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprMultiplicative) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
}

// /////////////////////////////////////

type ExprUnary struct {
	ASTNode

	Op    string      `parser:"[ @( OpBitwiseNot | OpLogicalNot | OpMinus | OpPlus ) ]" json:"op,omitempty"`
	Right ExprPostfix `parser:"@@"                                             json:"right"`
}

func (e *ExprUnary) Clone() *ExprUnary {
	if e == nil {
		return nil
	}

	return &ExprUnary{
		ASTNode: e.ASTNode.Clone(),
		Op:      e.Op,
		Right:   *e.Right.Clone(),
	}
}

func (e *ExprUnary) Children() (children []Node) {
	children = append(children, &e.Right)

	return
}

func (e ExprUnary) String() string {
	return fmt.Sprintf("%v%v", e.Op, e.Right)
}

// /////////////////////////////////////

type ExprPostfix struct {
	ASTNode

	Left  ExprPrimary `parser:"@@"             json:"left,omitempty"`
	Right *Expr       `parser:"[ '[' @@ ']' ]" json:"right,omitempty"`
}

func (e *ExprPostfix) Clone() *ExprPostfix {
	if e == nil {
		return nil
	}

	return &ExprPostfix{
		ASTNode: e.ASTNode.Clone(),
		Left:    *e.Left.Clone(),
		Right:   e.Right.Clone(),
	}
}

func (e *ExprPostfix) Children() (children []Node) {
	children = append(children, &e.Left)

	if e.Right != nil {
		children = append(children, e.Right)
	}

	return
}

func (e ExprPostfix) String() string {
	if e.Right != nil {
		return fmt.Sprintf("%v[%v]", e.Left, e.Right)
	}

	return e.Left.String()
}

// /////////////////////////////////////

type ExprPrimary struct {
	ASTNode

	SubExpression *Expr           `parser:"(  '(' @@ ')'  " json:"sub_expression,omitempty"`
	Invocation    *ExprInvocation `parser:"  | @@         " json:"invocation,omitempty"`
	Value         *Value          `parser:"  | @@        )" json:"value,omitempty"`
}

func (e *ExprPrimary) Clone() *ExprPrimary {
	if e == nil {
		return nil
	}

	return &ExprPrimary{
		ASTNode:       e.ASTNode.Clone(),
		SubExpression: e.SubExpression.Clone(),
		Invocation:    e.Invocation.Clone(),
		Value:         e.Value.Clone(),
	}
}

func (e *ExprPrimary) Children() (children []Node) {
	switch {
	case e.SubExpression != nil:
		children = append(children, e.SubExpression)
	case e.Invocation != nil:
		children = append(children, e.Invocation)
	case e.Value != nil:
		children = append(children, e.Value)

	}

	return
}

func (e ExprPrimary) String() string {
	switch {
	case e.SubExpression != nil:
		return e.SubExpression.String()
	case e.Invocation != nil:
		return e.Invocation.String()
	case e.Value != nil:
		return e.Value.String()
	default:
		return TokenNull
	}
}

// /////////////////////////////////////

type ExprInvocation struct {
	ASTNode

	Ident   Ident                   `parser:"@@"              json:"ident"`
	Monads  []*ExprInvocationParams `parser:"( '(' @@ ')' )+" json:"monads,omitempty"`
	Postfix *ExprPostfix            `parser:"[ '.' @@ ]"      json:"postfix,omitempty"`
}

func (e *ExprInvocation) Clone() *ExprInvocation {
	if e == nil {
		return nil
	}

	return &ExprInvocation{
		ASTNode: e.ASTNode.Clone(),
		Ident:   *e.Ident.Clone(),
		Monads:  cloneCollection(e.Monads),
		Postfix: e.Postfix.Clone(),
	}
}

func (e ExprInvocation) Children() (children []Node) {
	children = append(children, &e.Ident)

	for _, item := range e.Monads {
		children = append(children, item)
	}

	if e.Postfix != nil {
		children = append(children, e.Postfix)
	}

	return
}

func (e ExprInvocation) String() string {
	invocationParams := "()"
	if len(e.Monads) != 0 {
		params := make([]string, 0, len(e.Monads))
		for _, p := range e.Monads {
			params = append(params, fmt.Sprintf("(%v)", p))
		}
		invocationParams = strings.Join(params, "")
	}

	if e.Postfix != nil {
		return fmt.Sprintf("%v%v.%v", e.Ident, invocationParams, e.Postfix)
	}

	return fmt.Sprintf("%v%v", e.Ident, invocationParams)
}

// /////////////////////////////////////

type ExprInvocationParams struct {
	ASTNode

	Values []*Expr `parser:"[ @@ (',' @@)* ]" json:"values,omitempty"`
}

func (e *ExprInvocationParams) Clone() *ExprInvocationParams {
	if e == nil {
		return nil
	}

	return &ExprInvocationParams{
		ASTNode: e.ASTNode.Clone(),
		Values:  cloneCollection(e.Values),
	}
}

func (e *ExprInvocationParams) Children() (children []Node) {
	for _, item := range e.Values {
		children = append(children, item)
	}

	return
}

func (e ExprInvocationParams) String() string {
	params := make([]string, 0, len(e.Values))
	for _, p := range e.Values {
		params = append(params, p.String())
	}

	return strings.Join(params, ", ")
}
