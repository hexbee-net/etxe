package etx

import "fmt"

type Expr struct {
	Left *Conditional `parser:"@@"`
}

type Conditional struct {
	Condition    *LogicalOr `parser:"@@"`
	ConditionOp  string     `parser:"[ Whitespace? @OpCondition Whitespace?"`
	Left         *LogicalOr `parser:"  @@"`
	ConditionSep string     `parser:"  Whitespace? @OpColon Whitespace?"`
	Right        *LogicalOr `parser:"  @@ ]"`
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
	if e.ConditionOp != "" {
		return fmt.Sprintf("%s %s %s %s %s", e.Condition, e.ConditionOp, e.Left, e.ConditionSep, e.Right)
	}

	return e.Condition.String()
}

func (e *LogicalOr) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *LogicalAnd) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *BitwiseOr) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *BitwiseXor) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *BitwiseAnd) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *Equality) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *Relational) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *Shift) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *Additive) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *Multiplicative) String() string {
	switch {
	case e.Op != "":
		return fmt.Sprintf("%s %s %s", e.Left, e.Op, e.Right)
	case e.Left != nil:
		return e.Left.String()
	default:
		return ""
	}
}

func (e *Unary) String() string {
	if e.Postfix != nil {
		return e.Postfix.String()
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
	if e.SubExpression != nil {
		return e.SubExpression.String()
	}

	if e.Value != nil {
		return e.Value.String()
	}

	return ""
}
