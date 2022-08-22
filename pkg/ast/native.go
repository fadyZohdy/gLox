package ast

import "time"

type Clock struct{}

func (c *Clock) arity() int {
	return 0
}

func (c *Clock) call(i *Interpreter, a []any) any {
	return float64(time.Now().UnixNano() / int64(time.Millisecond))
}

func (c Clock) String() string {
	return "<native fn>"
}
