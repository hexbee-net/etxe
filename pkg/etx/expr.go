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
	Condition *ExprLogicalOr `parser:"If Whitespace @@ (Whitespace|NewLine)"`
	Left      *Expr          `parser:"BlockStart (Whitespace|NewLine)? @@? (Whitespace|NewLine)? BlockEnd (Whitespace|NewLine)?"`
	Right     *Expr          `parser:"[ Else (Whitespace|NewLine) BlockStart (Whitespace|NewLine)? @@ (Whitespace|NewLine)? BlockEnd ]"`
}

type ExprSwitch struct {
	Selector *ExprLogicalOr `parser:"Switch Whitespace @@ (Whitespace|NewLine) BlockStart (Whitespace|NewLine)?"`
	Cases    []*ExprCase    `parser:"@@* (Whitespace|NewLine)? BlockEnd"`
}

type ExprCase struct {
	Conditions []*ExprLogicalOr `parser:"(  (Whitespace|NewLine)? Case Whitespace @@ ( (Whitespace|NewLine)? ',' (Whitespace|NewLine)? @@ )* Whitespace? OpColon (Whitespace|NewLine)?"`
	Default    bool             `parser:" | (Whitespace|NewLine)? @'default'  Whitespace? OpColon (Whitespace|NewLine)? )"`
	Expr       *Expr            `parser:"BlockStart (Whitespace|NewLine)? @@ (Whitespace|NewLine)? BlockEnd (Whitespace|NewLine)?"`
}

// ExprConditional is a ternary expression.
//
// Ternaries are bad but necessary for Terraform compatibility, so they are
// included at the expression top level but `if` and `switch` just skip them
// and go straight to the next level.
type ExprConditional struct {
	Condition    *ExprLogicalOr `parser:"@@"`
	ConditionOp  string         `parser:"[ Whitespace? @OpCondition Whitespace?"`
	True         *ExprLogicalOr `parser:"  @@"`
	ConditionSep string         `parser:"  Whitespace? @OpColon Whitespace?"`
	False        *ExprLogicalOr `parser:"  @@ ]"`
}

type ExprLogicalOr struct {
	Left  *ExprLogicalAnd `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @OpLogicalOr Whitespace?"`
	Right *ExprLogicalOr  `parser:"  @@ ]"`
}

type ExprLogicalAnd struct {
	Left  *ExprBitwiseOr  `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @OpLogicalAnd Whitespace?"`
	Right *ExprLogicalAnd `parser:"  @@ ]"`
}

type ExprBitwiseOr struct {
	Left  *ExprBitwiseXor `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @OpBitwiseOr Whitespace?"`
	Right *ExprBitwiseOr  `parser:"  @@ ]"`
}

type ExprBitwiseXor struct {
	Left  *ExprBitwiseAnd `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @OpBitwiseXOr Whitespace?"`
	Right *ExprBitwiseXor `parser:"  @@ ]"`
}

type ExprBitwiseAnd struct {
	Left  *ExprEquality   `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @OpBitwiseAnd Whitespace?"`
	Right *ExprBitwiseAnd `parser:"  @@ ]"`
}

type ExprEquality struct {
	Left  *ExprRelational `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @( OpNotEqual | OpEqual ) Whitespace?"`
	Right *ExprEquality   `parser:"  @@ ]"`
}

type ExprRelational struct {
	Left  *ExprShift      `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @( OpMore | OpMoreOrEqual | OpLess | OpLessOrEqual ) Whitespace?"`
	Right *ExprRelational `parser:"  @@ ]"`
}

type ExprShift struct {
	Left  *ExprAdditive `parser:"@@"`
	Op    string        `parser:"[ Whitespace? @( OpBitwiseShiftLeft | OpBitwiseShiftRight ) Whitespace?"`
	Right *ExprShift    `parser:"  @@ ]"`
}

type ExprAdditive struct {
	Left  *ExprMultiplicative `parser:"@@"`
	Op    string              `parser:"[ Whitespace? @( OpMinus | OpPlus ) Whitespace?"`
	Right *ExprAdditive       `parser:"  @@ ]"`
}

type ExprMultiplicative struct {
	Left  *ExprUnary          `parser:"@@"`
	Op    string              `parser:"[ Whitespace? @( OpDivision | OpMultiplication | OpModulo ) Whitespace?"`
	Right *ExprMultiplicative `parser:"  @@ ]"`
}

type ExprUnary struct {
	Op      string       `parser:"  ( @( OpBitwiseNot | OpLogicalNot | OpMinus ) Whitespace?"`
	Unary   *ExprUnary   `parser:"    @@ )"`
	Postfix *ExprPostfix `parser:"| @@"`
}

type ExprPostfix struct {
	Left  *ExprPrimary `parser:"@@"`
	Right *Expr        `parser:"[ Whitespace? OpLBracket Whitespace? @@ Whitespace? OpRBracket Whitespace? ]"`
}

type ExprPrimary struct {
	SubExpression *Expr           `parser:"  OpLParen Whitespace? @@ Whitespace? OpRParen"`
	Invocation    *ExprInvocation `parser:"| @@"`
	Value         *Value          `parser:"| @@"`
}

type ExprInvocation struct {
	Ident      *Ident  `parser:"@@"`
	Parameters []*Expr `parser:"OpLParen (@@ ( Whitespace? ',' Whitespace? @@ )*)? OpRParen"`
}

// /////////////////////////////////////

func (e *Expr) String() string {
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

func (e *ExprIf) String() string {
	switch {
	case e.Condition == nil:
		panic("if condition cannot be <nil>")
	case e.Left == nil:
		return fmt.Sprintf("if %s { }", e.Condition)
	case e.Right == nil:
		return fmt.Sprintf("if %s {\n%s\n}", e.Condition, indent(e.Left.String(), indentationChar))
	default:
		return fmt.Sprintf("if %s {\n%s\n} else {\n%s\n}", e.Condition, indent(e.Left.String(), indentationChar), indent(e.Right.String(), indentationChar))

	}
}

func (e *ExprSwitch) String() string {
	if e.Selector == nil {
		panic("switch selector cannot be <nil>")
	}

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

func (e *ExprCase) String() string {
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

func (e *ExprConditional) String() string {
	switch {
	case e.ConditionOp != "" && e.ConditionSep == "",
		e.ConditionOp == "" && e.ConditionSep != "":
		panic("both operators need to be set")

	case e.ConditionOp != "" && e.ConditionSep != "":
		switch {
		case e.True == nil && e.False == nil:
			panic("true and false expressions must be set when operators are set")
		case e.True == nil:
			panic("true expression must be set when operators are set")
		case e.False == nil:
			panic("false expression must be set when operators are set")
		default:
			return fmt.Sprintf("%s %s %s %s %s", e.Condition, e.ConditionOp, e.True, e.ConditionSep, e.False)
		}

	case e.Condition != nil:
		return e.Condition.String()

	default:
		panic("condition cannot be <nil>")
	}
}

func (e *ExprLogicalOr) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprLogicalAnd) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprBitwiseOr) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprBitwiseXor) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprBitwiseAnd) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprEquality) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprRelational) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprShift) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprAdditive) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprMultiplicative) String() string {
	switch {
	case e.Op != "" && e.Right != nil:
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Op != "" && e.Right == nil:
		panic("operator with <nil> right side")
	case e.Left != nil:
		return e.Left.String()
	default:
		panic("left side cannot be <nil>")
	}
}

func (e *ExprUnary) String() string {
	if e.Postfix != nil {
		return e.Postfix.String()
	}

	if e.Unary == nil {
		panic("postfix and unary cannot both be <nil>")
	}

	return fmt.Sprintf("%s%s", e.Op, e.Unary)
}

func (e *ExprPostfix) String() string {
	if e.Right != nil {
		return fmt.Sprintf("%s[%s]", e.Left, e.Right)
	}

	return e.Left.String()
}

func (e *ExprPrimary) String() string {
	switch {
	case e.SubExpression != nil:
		return e.SubExpression.String()
	case e.Invocation != nil:
		return e.Invocation.String()
	case e.Value != nil:
		return e.Value.String()
	default:
		return ""
	}
}

func (e *ExprInvocation) String() string {
	params := make([]string, 0, len(e.Parameters))
	for _, p := range e.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("%s(%s)", e.Ident, strings.Join(params, ", "))
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
		Condition: e.Condition.Clone(),
		Left:      e.Left.Clone(),
		Right:     e.Right.Clone(),
	}
}

func (e *ExprSwitch) Clone() *ExprSwitch {
	if e == nil {
		return nil
	}

	out := &ExprSwitch{
		Selector: e.Selector.Clone(),
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
		out.Conditions = make([]*ExprLogicalOr, 0, len(e.Conditions))
		for _, c := range e.Conditions {
			out.Conditions = append(out.Conditions, c.Clone())
		}
	}

	return out
}

func (e *ExprConditional) Clone() *ExprConditional {
	if e == nil {
		return nil
	}

	return &ExprConditional{
		Condition:    e.Condition.Clone(),
		ConditionOp:  e.ConditionOp,
		True:         e.True.Clone(),
		ConditionSep: e.ConditionSep,
		False:        e.False.Clone(),
	}
}

func (e *ExprLogicalOr) Clone() *ExprLogicalOr {
	if e == nil {
		return nil
	}

	return &ExprLogicalOr{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprLogicalAnd) Clone() *ExprLogicalAnd {
	if e == nil {
		return nil
	}

	return &ExprLogicalAnd{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprBitwiseOr) Clone() *ExprBitwiseOr {
	if e == nil {
		return nil
	}

	return &ExprBitwiseOr{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprBitwiseXor) Clone() *ExprBitwiseXor {
	if e == nil {
		return nil
	}

	return &ExprBitwiseXor{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprBitwiseAnd) Clone() *ExprBitwiseAnd {
	if e == nil {
		return nil
	}

	return &ExprBitwiseAnd{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprEquality) Clone() *ExprEquality {
	if e == nil {
		return nil
	}

	return &ExprEquality{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprRelational) Clone() *ExprRelational {
	if e == nil {
		return nil
	}

	return &ExprRelational{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprShift) Clone() *ExprShift {
	if e == nil {
		return nil
	}

	return &ExprShift{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprAdditive) Clone() *ExprAdditive {
	if e == nil {
		return nil
	}

	return &ExprAdditive{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprMultiplicative) Clone() *ExprMultiplicative {
	if e == nil {
		return nil
	}

	return &ExprMultiplicative{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *ExprUnary) Clone() *ExprUnary {
	if e == nil {
		return nil
	}

	return &ExprUnary{
		Op:      e.Op,
		Unary:   e.Unary.Clone(),
		Postfix: e.Postfix.Clone(),
	}
}

func (e *ExprPostfix) Clone() *ExprPostfix {
	if e == nil {
		return nil
	}

	return &ExprPostfix{
		Left:  e.Left.Clone(),
		Right: e.Right.Clone(),
	}
}

func (e *ExprPrimary) Clone() *ExprPrimary {
	if e == nil {
		return nil
	}

	return &ExprPrimary{
		SubExpression: e.SubExpression.Clone(),
		Value:         e.Value.Clone(),
	}
}
