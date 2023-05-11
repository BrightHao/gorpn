package calc

import (
	"log"
	"strconv"
)

type iToken interface {
	Value() interface{} // 为了一个变量里存不同类型的 Token, 此处 Value 函数的返回值必须相同
}

// 操作数
type ValueToken struct {
	val float64
}

func (vt ValueToken) Value() interface{} {
	return vt.val
}

func (vt ValueToken) ValueString() string {
	return strconv.FormatFloat(vt.val, 'f', -1, 64)
}

func NewValueToken(s string) ValueToken {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal("字符串转为数值出错", err)
	}
	return ValueToken{
		val: f,
	}
}

type OpChar uint8

const (
	OpCharPlus     OpChar = '+'
	OpCharMinus    OpChar = '-'
	OpCharMultiple OpChar = '*'
	OpCharDivide   OpChar = '/'
	OpCharPow      OpChar = '^'
	OpCharAnd      OpChar = '&'
	OpCharOr       OpChar = '|'
	OpCharNot      OpChar = '!'
)

// 操作符类型
type OpType int

const (
	OpTypePlus     OpType = iota // 加
	OpTypeMinus                  // 减
	OpTypeMultiple               // 乘
	OpTypeDivide                 // 除
	OpTypePow                    // 幂
	OpTypeAnd                    // 与
	OpTypeOr                     // 或
)

// 操作符优先级, 值越小, 优先级越高, 相等的优先级一样
type OpPriority int

const (
	OpPriorPlus     OpPriority = 100
	OpPriorMinus    OpPriority = 100
	OpPriorMultiple OpPriority = 80
	OpPriorDivide   OpPriority = 80
	OpPriorPow      OpPriority = 60
	OpPriorAnd      OpPriority = 20
	OpPriorOr       OpPriority = 20
)

// 操作符
type OpToken struct {
	typ   OpType
	prior OpPriority
	val   OpChar
}

func (ot OpToken) Value() interface{} {
	return ot.val
}

func (ot OpToken) Prior() OpPriority {
	return ot.prior
}

func NewPlusOperator() OpToken {
	return OpToken{
		val:   OpCharPlus,
		typ:   OpTypePlus,
		prior: OpPriorPlus,
	}
}

func NewMinusOperator() OpToken {
	return OpToken{
		val:   OpCharMinus,
		typ:   OpTypeMinus,
		prior: OpPriorMinus,
	}
}

func NewMultipleOperator() OpToken {
	return OpToken{
		val:   OpCharMultiple,
		typ:   OpTypeMultiple,
		prior: OpPriorMultiple,
	}
}

func NewDivideOperator() OpToken {
	return OpToken{
		val:   OpCharDivide,
		typ:   OpTypeDivide,
		prior: OpPriorDivide,
	}
}

func NewPowOperator() OpToken {
	return OpToken{
		val:   OpCharPow,
		typ:   OpTypePow,
		prior: OpPriorPow,
	}
}

func NewAndOperator() OpToken {
	return OpToken{
		val:   OpCharAnd,
		typ:   OpTypeAnd,
		prior: OpPriorAnd,
	}
}

func NewOrOperator() OpToken {
	return OpToken{
		val:   OpCharOr,
		typ:   OpTypeOr,
		prior: OpPriorOr,
	}
}

type BracketChar uint8

const (
	BracketCharLeft  BracketChar = '('
	BracketCharRight BracketChar = ')'
)

type BracketType int

const (
	BracketTypeLeft BracketType = iota
	BracketTypeRight
)

type BracketToken struct {
	typ BracketType
	val BracketChar
}

func (bt BracketToken) Value() interface{} {
	return bt.val
}

func NewLeftBracketToken() BracketToken {
	return BracketToken{
		val: BracketCharLeft,
		typ: BracketTypeLeft,
	}
}

func NewRightBracketToken() BracketToken {
	return BracketToken{
		val: BracketCharRight,
		typ: BracketTypeRight,
	}
}
