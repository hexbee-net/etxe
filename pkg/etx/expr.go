package etx

import (
	"fmt"
	"strings"
)

type Expr struct {
	Left   *ExprConditional `parser:"(   @@  " json:"left,omitempty"`
	If     *ExprIf          `parser:"  | @@  " json:"if,omitempty"`
	Switch *ExprSwitch      `parser:"  | @@ )" json:"switch,omitempty"`
}

type ExprIf struct {
	Condition ExprLogicalOr `parser:"If @@ NewLine? '{' NewLine?"           json:"condition"`
	Left      *Expr         `parser:"@@? NewLine? '}' NewLine?"             json:"left,omitempty"`
	Right     *Expr         `parser:"[ Else '{' NewLine? @@ NewLine? '}' ]" json:"right,omitempty"`
}

type ExprSwitch struct {
	Selector ExprLogicalOr `parser:"Switch @@ NewLine? '{' NewLine?"   json:"selector"`
	Cases    []*ExprCase   `parser:"(@@ NewLine?)* '}'"                json:"cases,omitempty"`
}

type ExprCase struct {
	Conditions []ExprLogicalOr `parser:"(   Case @@ ( ',' @@ )* ':'  "         json:"conditions"`
	Default    bool            `parser:"  | @'default'          ':' )"         json:"default"`
	Expr       *Expr           `parser:"NewLine? '{' NewLine? @@ NewLine? '}'" json:"expr,omitempty"`
}

// ExprConditional is a ternary expression.
//
// Ternaries are bad but necessary for Terraform compatibility, so they are
// included at the expression top level but `if` and `switch` just skip them
// and go straight to the next level.
type ExprConditional struct {
	Condition   ExprLogicalOr  `parser:"@@"          json:"condition"`
	ConditionOp bool           `parser:"[ @'?'     " json:"condition_op,omitempty"`
	TrueExpr    *ExprLogicalOr `parser:"  @@       " json:"true_expr,omitempty"`
	FalseExpr   *ExprLogicalOr `parser:"  ':' @@  ]" json:"false_expr,omitempty"`
}

type ExprLogicalOr struct {
	Left  ExprLogicalAnd `parser:"@@"               json:"left"`
	Op    string         `parser:"[ @OpLogicalOr  " json:"op,omitempty"`
	Right *ExprLogicalOr `parser:"  @@           ]" json:"right,omitempty"`
}

type ExprLogicalAnd struct {
	Left  ExprBitwiseOr   `parser:"@@"                json:"left"`
	Op    string          `parser:"[ @OpLogicalAnd  " json:"op,omitempty"`
	Right *ExprLogicalAnd `parser:"  @@            ]" json:"right,omitempty"`
}

type ExprBitwiseOr struct {
	Left  ExprBitwiseXor `parser:"@@"               json:"left"`
	Op    string         `parser:"[ @OpBitwiseOr  " json:"op,omitempty"`
	Right *ExprBitwiseOr `parser:"  @@           ]" json:"right,omitempty"`
}

type ExprBitwiseXor struct {
	Left  ExprBitwiseAnd  `parser:"@@"                json:"left"`
	Op    string          `parser:"[ @OpBitwiseXOr  " json:"op,omitempty"`
	Right *ExprBitwiseXor `parser:"  @@            ]" json:"right,omitempty"`
}

type ExprBitwiseAnd struct {
	Left  ExprEquality    `parser:"@@"                json:"left"`
	Op    string          `parser:"[ @OpBitwiseAnd  " json:"op,omitempty"`
	Right *ExprBitwiseAnd `parser:"  @@            ]" json:"right,omitempty"`
}

type ExprEquality struct {
	Left  ExprRelational `parser:"@@"                            json:"left"`
	Op    string         `parser:"[ @( OpNotEqual | OpEqual )  " json:"op,omitempty"`
	Right *ExprEquality  `parser:"  @@                        ]" json:"right,omitempty"`
}

type ExprRelational struct {
	Left  ExprShift       `parser:"@@"                                                       json:"left"`
	Op    string          `parser:"[ @( OpMore | OpMoreOrEqual | OpLess | OpLessOrEqual )  " json:"op,omitempty"`
	Right *ExprRelational `parser:"  @@                                                   ]" json:"right,omitempty"`
}

type ExprShift struct {
	Left  ExprAdditive `parser:"@@"                                                json:"left"`
	Op    string       `parser:"[ @( OpBitwiseShiftLeft | OpBitwiseShiftRight )  " json:"op,omitempty"`
	Right *ExprShift   `parser:"  @@                                            ]" json:"right,omitempty"`
}

type ExprAdditive struct {
	Left  ExprMultiplicative `parser:"@@"                        json:"left"`
	Op    string             `parser:"[ @( OpMinus | OpPlus )  " json:"op,omitempty"`
	Right *ExprAdditive      `parser:"  @@                    ]" json:"right,omitempty"`
}

type ExprMultiplicative struct {
	Left  ExprUnary           `parser:"@@"                                                json:"left,omitempty"`
	Op    string              `parser:"[ @( OpDivision | OpMultiplication | OpModulo )  " json:"op,omitempty"`
	Right *ExprMultiplicative `parser:"  @@                                            ]" json:"right,omitempty"`
}

type ExprUnary struct {
	Op    string      `parser:"[ @( OpBitwiseNot | OpLogicalNot | OpMinus ) ]" json:"op,omitempty"`
	Right ExprPostfix `parser:"@@"                                             json:"right"`
}

type ExprPostfix struct {
	Left  ExprPrimary `parser:"@@"             json:"left,omitempty"`
	Right *Expr       `parser:"[ '[' @@ ']' ]" json:"right,omitempty"`
}

type ExprPrimary struct {
	SubExpression *Expr           `parser:"(  '(' @@ ')'  " json:"sub_expression,omitempty"`
	Invocation    *ExprInvocation `parser:"  | @@         " json:"invocation,omitempty"`
	Value         *Value          `parser:"  | @@        )" json:"value,omitempty"`
}

type ExprInvocation struct {
	Ident   Ident                  `parser:"@@"              json:"ident"`
	Monads  []ExprInvocationParams `parser:"( '(' @@ ')' )+" json:"monads,omitempty"`
	Postfix *ExprPostfix           `parser:"[ '.' @@ ]"      json:"postfix,omitempty"`
}

type ExprInvocationParams struct {
	Values []Expr `parser:"[ @@ (',' @@)* ]" json:"values,omitempty"`
}

// /////////////////////////////////////

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

func (e ExprLogicalOr) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
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

func (e ExprBitwiseOr) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
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

func (e ExprBitwiseAnd) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
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

func (e ExprRelational) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
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

func (e ExprAdditive) String() string {
	if e.Op == "" {
		return e.Left.String()
	}

	if e.Right == nil {
		panic("operator with <nil> right side")
	}

	return fmt.Sprintf("%v %v %v", e.Left, e.Op, e.Right)
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

func (e ExprUnary) String() string {
	return fmt.Sprintf("%v%v", e.Op, e.Right)
}

func (e ExprPostfix) String() string {
	if e.Right != nil {
		return fmt.Sprintf("%v[%v]", e.Left, e.Right)
	}

	return e.Left.String()
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

func (e ExprInvocationParams) String() string {
	params := make([]string, 0, len(e.Values))
	for _, p := range e.Values {
		params = append(params, p.String())
	}

	return strings.Join(params, ", ")
}

// /////////////////////////////////////

func (e *Expr) Clone() *Expr {
	if e == nil {
		return nil
	}

	return &Expr{
		Left:   e.Left.Clone(),
		If:     e.If.Clone(),
		Switch: e.Switch.Clone(),
	}
}

func (e *ExprIf) Clone() *ExprIf {
	if e == nil {
		return nil
	}

	return &ExprIf{
		Condition: *e.Condition.Clone(),
		Left:      e.Left.Clone(),
		Right:     e.Right.Clone(),
	}
}

func (e *ExprSwitch) Clone() *ExprSwitch {
	if e == nil {
		return nil
	}

	out := &ExprSwitch{
		Selector: *e.Selector.Clone(),
		Cases:    nil,
	}

	if e.Cases != nil {
		out.Cases = make([]*ExprCase, 0, len(e.Cases))
		for _, c := range e.Cases {
			out.Cases = append(out.Cases, c.Clone())
		}
	}

	return out
}

func (e *ExprCase) Clone() *ExprCase {
	if e == nil {
		return nil
	}

	out := &ExprCase{
		Conditions: nil,
		Default:    e.Default,
		Expr:       e.Expr.Clone(),
	}

	if e.Conditions != nil {
		out.Conditions = make([]ExprLogicalOr, 0, len(e.Conditions))
		for _, c := range e.Conditions {
			out.Conditions = append(out.Conditions, *c.Clone())
		}
	}

	return out
}

func (e *ExprConditional) Clone() *ExprConditional {
	if e == nil {
		return nil
	}

	return &ExprConditional{
		Condition:   *e.Condition.Clone(),
		ConditionOp: e.ConditionOp,
		TrueExpr:    e.TrueExpr.Clone(),
		FalseExpr:   e.FalseExpr.Clone(),
	}
}

func (e *ExprLogicalOr) Clone() *ExprLogicalOr {
	if e == nil {
		return nil
	}

	return &ExprLogicalOr{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprLogicalAnd) Clone() *ExprLogicalAnd {
	if e == nil {
		return nil
	}

	return &ExprLogicalAnd{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprBitwiseOr) Clone() *ExprBitwiseOr {
	if e == nil {
		return nil
	}

	return &ExprBitwiseOr{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprBitwiseXor) Clone() *ExprBitwiseXor {
	if e == nil {
		return nil
	}

	return &ExprBitwiseXor{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprBitwiseAnd) Clone() *ExprBitwiseAnd {
	if e == nil {
		return nil
	}

	return &ExprBitwiseAnd{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprEquality) Clone() *ExprEquality {
	if e == nil {
		return nil
	}

	return &ExprEquality{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprRelational) Clone() *ExprRelational {
	if e == nil {
		return nil
	}

	return &ExprRelational{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprShift) Clone() *ExprShift {
	if e == nil {
		return nil
	}

	return &ExprShift{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprAdditive) Clone() *ExprAdditive {
	if e == nil {
		return nil
	}

	return &ExprAdditive{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprMultiplicative) Clone() *ExprMultiplicative {
	if e == nil {
		return nil
	}

	return &ExprMultiplicative{
		Left:  *e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprUnary) Clone() *ExprUnary {
	if e == nil {
		return nil
	}

	return &ExprUnary{
		Op:    e.Op,
		Right: *e.Right.Clone(),
	}
}

func (e *ExprPostfix) Clone() *ExprPostfix {
	if e == nil {
		return nil
	}

	return &ExprPostfix{
		Left:  *e.Left.Clone(),
		Right: e.Right.Clone(),
	}
}

func (e *ExprPrimary) Clone() *ExprPrimary {
	if e == nil {
		return nil
	}

	return &ExprPrimary{
		SubExpression: e.SubExpression.Clone(),
		Invocation:    e.Invocation.Clone(),
		Value:         e.Value.Clone(),
	}
}

func (e *ExprInvocation) Clone() *ExprInvocation {
	if e == nil {
		return nil
	}

	out := &ExprInvocation{
		Ident:   *e.Ident.Clone(),
		Monads:  nil,
		Postfix: e.Postfix.Clone(),
	}

	if e.Monads == nil {
		return out
	}

	out.Monads = make([]ExprInvocationParams, 0, len(e.Monads))
	for _, p := range e.Monads {
		out.Monads = append(out.Monads, *p.Clone())
	}

	return out
}

func (e *ExprInvocationParams) Clone() *ExprInvocationParams {
	if e == nil {
		return nil
	}

	out := &ExprInvocationParams{
		Values: nil,
	}

	if e.Values == nil {
		return out
	}

	out.Values = make([]Expr, 0, len(e.Values))
	for _, p := range e.Values {
		out.Values = append(out.Values, *p.Clone())
	}

	return out
}
