package calc

import (
	"log"
	"math"
	"strings"

	"github.com/BrightHao/gorpn/stack"
	"github.com/BrightHao/gorpn/util"
)

// 计算器对象，用逆波兰法求值
// 算法可以参考 https://zhuanlan.zhihu.com/p/40638139
type Calc struct {
	s  stack.Stack // 求逆波兰表达式的栈
	s1 stack.Stack // 运算符栈
	s2 stack.Stack // 操作数栈
	rp string      // 逆波兰表达式字符串

	parser Parser // 字符串解析器
}

func NewCalc(exp string) (c Calc) {
	p := Parser{exp: exp}
	p.exp2Tokens()

	c = Calc{
		s:      stack.NewStack(),
		s1:     stack.NewStack(),
		s2:     stack.NewStack(),
		parser: p,
	}
	return
}

// 转换逆波兰表达式
func (cal *Calc) Result() (res float64) {
	// 1. 初始化两个栈：运算符栈s1和储存中间结果的栈s2；
	// 2. 从左至右扫描中缀表达式；
	for _, token := range cal.parser.tokens {
		switch token.(type) {
		case ValueToken:
			// 3. 遇到操作数时，将其压s2；
			cal.s2.Push(token)
		case OpToken:
			t, _ := token.(OpToken)
			// 4. 遇到运算符时，比较其与s1栈顶运算符的优先级：
		check4_1:
			if tt, ok := cal.s1.Peek().(BracketToken); ok && tt.typ == BracketTypeLeft || cal.s1.Len() == 0 {
				// 4.1. 如果s1为空，或栈顶运算符为左括号“(”，则直接将此运算符入栈；
				cal.s1.Push(token)
			} else if tt, ok := cal.s1.Peek().(OpToken); ok {
				if t.prior < tt.prior {
					// 4.2. 否则，若优先级比栈顶运算符的高，也将运算符压入s1（注意转换为前缀表达式时是优先级较高或相同，而这里则不包括相同的情况）；
					cal.s1.Push(token)
				} else {
					// 4.3. 否则，将s1栈顶的运算符弹出并压入到s2中，再次转到(4.1)与s1中新的栈顶运算符相比较；
					cal.s2.Push(cal.s1.Pop())
					goto check4_1
				}
			} else if !ok {
				log.Fatal("希望遇到运算符, 但是没有")
			}
		case BracketToken:
			// 5. 遇到括号时：
			// 5.1. 如果是左括号“(”，则直接压入s1；
			if token.(BracketToken).typ == BracketTypeLeft {
				cal.s1.Push(token)
			}
			// 5.2. 如果是右括号“)”，则依次弹出s1栈顶的运算符，并压入s2，直到遇到左括号为止，此时将这一对括号丢弃；
			if token.(BracketToken).typ == BracketTypeRight {
				for {
					t := cal.s1.Pop()
					if tt, ok := t.(OpToken); ok {
						cal.s2.Push(tt)
					} else if tt, ok := t.(BracketToken); ok {
						if tt.typ == BracketTypeLeft {
							// 遇到左括号停止
							break
						} else {
							log.Fatal("遇到的是括号, 但是不是左括号, 说明表达式的逻辑有问题")
						}
					} else {
						log.Fatal("遇到了运算符和括号之外的东西, 处理失败")
					}
				}
			}
		default:
		}
	}
	// 6. 直到表达式的最右边；
	// 7. 将s1中剩余的运算符依次弹出并压入s2；
	for {
		if cal.s1.Len() == 0 {
			break
		}
		t := cal.s1.Pop()
		if _, ok := t.(OpToken); ok {
			cal.s2.Push(t)
		} else {
			log.Fatal("栈 s1 中剩余非运算符, 很可能表达式不合法")
		}
	}
	// 8. 依次弹出s2中的元素并输出，结果的逆序即为中缀表达式对应的后缀表达式（转换为前缀表达式时不用逆序）
	var rpSlice []string       // 倒序的 逆波兰表达式 的切片
	var rpTokens []interface{} // 倒序的 tokens 的切片
	for {
		if cal.s2.Len() == 0 {
			break
		}
		t := cal.s2.Pop()

		rpTokens = append(rpTokens, t)

		switch t.(type) {
		case ValueToken:
			rpSlice = append(rpSlice, t.(ValueToken).ValueString())
		case OpToken:
			rpSlice = append(rpSlice, string(t.(OpToken).Value().(OpChar)))
		default:
			log.Fatal("s2 弹出的时候遇到未知类型的 token")
		}
	}
	util.ReverseAny(rpSlice) // 翻转成从做到右的顺序
	cal.rp = strings.Join(rpSlice, " ")

	util.ReverseAny(rpTokens)
	// 循环从左到右正序的 tokens 切片, 借助 s 栈求表达式的值
	// 从左至右扫描表达式，
	// 遇到数字时，将数字压入堆栈，
	// 遇到运算符时，弹出栈顶的两个数，用运算符对它们做相应的计算（次顶元素 op 栈顶元素），并将结果入栈；
	// 重复上述过程直到表达式最右端，最后运算得出的值即为表达式的结果
	for _, t := range rpTokens {
		switch t.(type) {
		case ValueToken:
			cal.s.Push(t.(ValueToken).Value().(float64))
		case OpToken:
			n1 := cal.s.Pop().(float64)
			n2 := cal.s.Pop().(float64)
			switch t.(OpToken).typ {
			case OpTypePlus:
				cal.s.Push(n2 + n1)
			case OpTypeMinus:
				cal.s.Push(n2 - n1)
			case OpTypeMultiple:
				cal.s.Push(n2 * n1)
			case OpTypeDivide:
				cal.s.Push(n2 / n1)
			case OpTypePow:
				cal.s.Push(math.Pow(n2, n1))
			case OpTypeAnd:
				cal.s.Push(float64(int64(n2) & int64(n1)))
			case OpTypeOr:
				cal.s.Push(float64(int64(n2) | int64(n1)))
			default:
				log.Fatal("不支持的运算符", string(t.(OpToken).Value().(OpChar)))
			}
		default:
			log.Fatal("逆波兰表达式求值时遇到数值和运算符之外的 token")
		}
	}
	res = cal.s.Pop().(float64)

	return
}
