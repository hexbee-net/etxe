package etx

import "fmt"

type Expr struct {
	Left *Conditional `parser:"@@"`
}

type Conditional struct {
	Condition    *LogicalOr `parser:"@@"`
	ConditionOp  string     `parser:"[ Whitespace? @OpCondition Whitespace?"`
	True         *LogicalOr `parser:"  @@"`
	ConditionSep string     `parser:"  Whitespace? @OpColon Whitespace?"`
	False        *LogicalOr `parser:"  @@ ]"`
}

type LogicalOr struct {
	Left  *LogicalAnd `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @OpLogicalOr Whitespace?"`
	Right *LogicalOr  `parser:"  @@ ]"`
}

type LogicalAnd struct {
	Left  *BitwiseOr  `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @OpLogicalAnd Whitespace?"`
	Right *LogicalAnd `parser:"  @@ ]"`
}

type BitwiseOr struct {
	Left  *BitwiseXor `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @OpBitwiseOr Whitespace?"`
	Right *BitwiseOr  `parser:"  @@ ]"`
}

type BitwiseXor struct {
	Left  *BitwiseAnd `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @OpBitwiseXOr Whitespace?"`
	Right *BitwiseXor `parser:"  @@ ]"`
}

type BitwiseAnd struct {
	Left  *Equality   `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @OpBitwiseAnd Whitespace?"`
	Right *BitwiseAnd `parser:"  @@ ]"`
}

type Equality struct {
	Left  *Relational `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @( OpNotEqual | OpEqual ) Whitespace?"`
	Right *Equality   `parser:"  @@ ]"`
}

type Relational struct {
	Left  *Shift      `parser:"@@"`
	Op    string      `parser:"[ Whitespace? @( OpMore | OpMoreOrEqual | OpLess | OpLessOrEqual ) Whitespace?"`
	Right *Relational `parser:"  @@ ]"`
}

type Shift struct {
	Left  *Additive `parser:"@@"`
	Op    string    `parser:"[ Whitespace? @( OpBitwiseShiftLeft | OpBitwiseShiftRight ) Whitespace?"`
	Right *Shift    `parser:"  @@ ]"`
}

type Additive struct {
	Left  *Multiplicative `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @( OpMinus | OpPlus ) Whitespace?"`
	Right *Additive       `parser:"  @@ ]"`
}

type Multiplicative struct {
	Left  *Unary          `parser:"@@"`
	Op    string          `parser:"[ Whitespace? @( OpDivision | OpMultiplication | OpModulo ) Whitespace?"`
	Right *Multiplicative `parser:"  @@ ]"`
}

type Unary struct {
	Op      string   `parser:"  ( @( OpBitwiseNot | OpLogicalNot | OpMinus ) Whitespace?"`
	Unary   *Unary   `parser:"    @@ )"`
	Postfix *Postfix `parser:"| @@"`
}

type Postfix struct {
	Left  *Primary `parser:"@@"`
	Right *Expr    `parser:"[ Whitespace? OpLBracket Whitespace? @@ Whitespace? OpRBracket Whitespace? ]"`
}

type Primary struct {
	SubExpression *Expr  `parser:"  OpLParen Whitespace? @@ Whitespace? OpRParen"`
	Value         *Value `parser:"| @@"`
}

// /////////////////////////////////////

func (e *Expr) String() string {
	return e.Left.String()
}

func (e *Conditional) String() string {
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

func (e *LogicalOr) String() string {
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

func (e *LogicalAnd) String() string {
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

func (e *BitwiseOr) String() string {
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

func (e *BitwiseXor) String() string {
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

func (e *BitwiseAnd) String() string {
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

func (e *Equality) String() string {
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

func (e *Relational) String() string {
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

func (e *Shift) String() string {
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

func (e *Additive) String() string {
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

func (e *Multiplicative) String() string {
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

func (e *Unary) String() string {
	if e.Postfix != nil {
		return e.Postfix.String()
	}

	if e.Unary == nil {
		panic("postfix and unary cannot both be <nil>")
	}

	return fmt.Sprintf("%s%s", e.Op, e.Unary)
}

func (e *Postfix) String() string {
	if e.Right != nil {
		return fmt.Sprintf("%s[%s]", e.Left, e.Right)
	}

	return e.Left.String()
}

func (e *Primary) String() string {
	switch {
	case e.SubExpression != nil:
		return e.SubExpression.String()
	case e.Value != nil:
		return e.Value.String()
	default:
		return ""
	}
}

// /////////////////////////////////////

func (e *Expr) Clone() *Expr {
	if e == nil {
		return nil
	}

	return &Expr{
		Left: e.Left.Clone(),
	}
}

func (e *Conditional) Clone() *Conditional {
	if e == nil {
		return nil
	}

	return &Conditional{
		Condition:    e.Condition.Clone(),
		ConditionOp:  e.ConditionOp,
		True:         e.True.Clone(),
		ConditionSep: e.ConditionSep,
		False:        e.False.Clone(),
	}
}

func (e *LogicalOr) Clone() *LogicalOr {
	if e == nil {
		return nil
	}

	return &LogicalOr{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *LogicalAnd) Clone() *LogicalAnd {
	if e == nil {
		return nil
	}

	return &LogicalAnd{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *BitwiseOr) Clone() *BitwiseOr {
	if e == nil {
		return nil
	}

	return &BitwiseOr{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *BitwiseXor) Clone() *BitwiseXor {
	if e == nil {
		return nil
	}

	return &BitwiseXor{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *BitwiseAnd) Clone() *BitwiseAnd {
	if e == nil {
		return nil
	}

	return &BitwiseAnd{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *Equality) Clone() *Equality {
	if e == nil {
		return nil
	}

	return &Equality{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *Relational) Clone() *Relational {
	if e == nil {
		return nil
	}

	return &Relational{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *Shift) Clone() *Shift {
	if e == nil {
		return nil
	}

	return &Shift{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *Additive) Clone() *Additive {
	if e == nil {
		return nil
	}

	return &Additive{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *Multiplicative) Clone() *Multiplicative {
	if e == nil {
		return nil
	}

	return &Multiplicative{
		Left:  e.Left.Clone(),
		Op:    e.Op,
		Right: e.Right.Clone(),
	}
}

func (e *Unary) Clone() *Unary {
	if e == nil {
		return nil
	}

	return &Unary{
		Op:      e.Op,
		Unary:   e.Unary.Clone(),
		Postfix: e.Postfix.Clone(),
	}
}

func (e *Postfix) Clone() *Postfix {
	if e == nil {
		return nil
	}

	return &Postfix{
		Left:  e.Left.Clone(),
		Right: e.Right.Clone(),
	}
}

func (e *Primary) Clone() *Primary {
	if e == nil {
		return nil
	}

	return &Primary{
		SubExpression: e.SubExpression.Clone(),
		Value:         e.Value.Clone(),
	}
}
